package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"users-api/database"
	"users-api/models"
	"users-api/repository"
)

func Signup(ctx *gin.Context) {
	userRepo := repository.NewUserRepository(database.DB)
	apiKey := ctx.Request.Header.Get("X-API-KEY")
	envApiKey := os.Getenv("X-API-KEY")

	if apiKey != envApiKey || apiKey == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Status   int16  `json:"status"`
		RoleId   int    `json:"role_id"`
	}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	body.Password = string(hash)

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		Status:   body.Status,
		RoleID:   body.RoleId,
	}
	if err := userRepo.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}
