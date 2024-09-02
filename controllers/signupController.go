package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"users-api/database"
	"users-api/models"
)

func Signup(ctx *gin.Context) {
	apiKey := ctx.Request.Header.Get("X-API-KEY")
	envApiKey := os.Getenv("X-API-KEY")

	if apiKey != envApiKey || apiKey == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = string(hash)
	if err := database.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}
