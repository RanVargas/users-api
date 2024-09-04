package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

func Login(ctx *gin.Context) {

	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBind(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := userRepo.GetUserByEmail(data.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Email or Password"})
		} else {
			log.Printf("Database error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}
	if data.Email != user.Email {
		log.Fatal("Mismatch between emails")
	}

	if user.ID == 0 {
		log.Printf("User Found but ID is 0 from email: %s", data.Email)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}
	userPassword, _ := userPasswordRepo.GetUserPassword(user.ID)

	err = bcrypt.CompareHashAndPassword([]byte(userPassword.UserPassword), []byte(data.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	log.Printf("Creating token for user ID: %d", user.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}
	log.Printf("Created token: %s", tokenString)
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"Authorization": tokenString,
	})
	//ctx.SetSameSite(http.SameSiteLaxMode)
	//ctx.SetCookie("Authorization", tokenString, 3600*2, "/", "", false, true)
}

func Validate(ctx *gin.Context) {
	_, err := ctx.Get("user")
	if err != false {
		ctx.JSON(500, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": "You are logged in"})
}

func Logout(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"data": "You are logged out"})
}
