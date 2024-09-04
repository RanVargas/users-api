package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"users-api/models"
)

func Signup(ctx *gin.Context) {

	apiKey := ctx.Request.Header.Get("X-API-KEY")
	envApiKey := os.Getenv("X-API-KEY")

	if apiKey != envApiKey || apiKey == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	var body struct {
		Name      string   `json:"name" binding:"required"`
		Email     string   `json:"email" binding:"required"`
		Password  string   `json:"password" binding:"required"`
		Status    int16    `json:"status" binding:"required"`
		RoleUid   string   `json:"role_uid" binding:"required"`
		GroupsUid []string `json:"groups_uid" binding:"required"`
	}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not process request body"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	body.Password = string(hash)
	role, err := roleRepo.FindRoleByUid(body.RoleUid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Provided role not found."})
		return
	}
	groups := make([]models.Group, len(body.GroupsUid))
	for i := 0; i < len(body.GroupsUid); i++ {
		group, e := groupRepo.FindGroupByUid(body.GroupsUid[i])
		if e != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Provided group not found."})
			return
		}
		groups = append(groups, *group)
	}

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		Status:   body.Status,
		RoleID:   role.ID,
		Group:    groups,
	}

	if err = userRepo.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user, details: " + err.Error()})
		return
	}
	if user.ID == 0 {
		log.Fatalf("Could not get ID")
	}

	userPassword := models.UserPassword{
		UserID:       user.ID,
		UserPassword: user.Password,
	}
	if err = userPasswordRepo.CreateUserPassword(&userPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user-password model, details:  " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}
