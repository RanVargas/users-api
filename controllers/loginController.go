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
	log.Printf("This is the user found: %v", user)

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 6).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"Authorization": tokenString,
	})

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
	ctx.Set("Authorization", "")
	ctx.JSON(http.StatusOK, gin.H{"result": "You are logged out"})
}
