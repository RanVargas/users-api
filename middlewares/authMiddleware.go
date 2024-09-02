package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"users-api/database"
	"users-api/models"
)

func RequireAuth(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		log.Printf("No authorization token was found and thus returning")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenString = strings.TrimSpace(tokenString)
	workedTokenString := strings.TrimPrefix(tokenString, "Bearer ")

	log.Printf("This is the token gotten: %s", workedTokenString)
	parts := strings.Split(workedTokenString, ".")
	if len(parts) != 3 {
		log.Printf("Token does not have the correct structure, expected 3 parts but got %d", len(parts))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(workedTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Println("JWT_SECRET environment variable is not set")
			return nil, fmt.Errorf("JWT secret key is missing")
		}
		return []byte(secret), nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while parsing JWT token"})
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.User
		userId := int(claims["id"].(float64))
		log.Printf("This was the token received from middleware: %s", tokenString)
		log.Printf("This is the id found in the JWT %v", userId)
		if err := database.DB.First(&user, userId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println(fmt.Sprintf("The id gotten from database is: %s, as such it has not been possible to continue", user.Id))
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Id retrieval error"})
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		if user.Id == 0 {
			fmt.Println("Id has come as 0")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", user)
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	ctx.Next()
}
