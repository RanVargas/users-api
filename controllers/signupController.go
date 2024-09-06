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

	type groupBody struct {
		Name  string `json:"name"`
		Uid   string `json:"Uid"`
		Users string `json:"-"`
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
	if len(body.GroupsUid) < 0 {
		log.Printf("Group could not be parsed correctly")
	}
	var groups []models.Group
	for _, groupUID := range body.GroupsUid {
		if groupUID == "" {
			continue
		}
		group, _ := groupRepo.FindGroupByUidWithNoUsers(groupUID)
		if group != nil {
			groups = append(groups, *group)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Provided group not found."})
			return
		}
	}
	if groups[0].Uid.String() == "" {
		log.Printf("Group could not be parsed correctly")
	}
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		Status:   body.Status,
		RoleID:   role.ID,
	}

	if err = userRepo.CreateUser(&user, groups); err != nil {
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
