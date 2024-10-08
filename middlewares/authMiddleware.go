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
	"users-api/repository"
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

		userRepo := repository.NewUserRepository(database.DB)
		userId := int(claims["id"].(float64))
		log.Printf("This is the id gotten from the conversion: %v", userId)
		if userId == 0 {
			log.Printf("User with id %v was not found", claims["id"].(float64))
		}
		user, er := userRepo.GetUserById(userId)
		if er != nil {
			if errors.Is(er, gorm.ErrRecordNotFound) {
				fmt.Println(fmt.Sprintf("The id gotten from database is: %s, as such it has not been possible to continue", user.ID))
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Id retrieval error"})
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			log.Printf("Error getting user: %v", er)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Id retrieval error"})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if user.ID == 0 {
			log.Printf("Id has come as 0")
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
