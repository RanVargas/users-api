package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"users-api/database"
	"users-api/models"
	"users-api/repository"
)

var roleRepo *repository.RoleRepository

func InitializeRolesRepo() {
	roleRepo = repository.NewRoleRepository(database.DB)
}

func CreateRole(ctx *gin.Context) {
	var role models.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := database.DB.Create(&role)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	}

	ctx.JSON(http.StatusCreated, role)
}

func GetAllRoles(ctx *gin.Context) {
	var roles []models.Role
	result := database.DB.Find(&roles)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, roles)
}

func UpdateRole(ctx *gin.Context) {
	var role models.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := roleRepo.UpdateRole(role); err != nil {
		log.Printf("Error updating role: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error updating user"})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func DeleteRole(ctx *gin.Context) {
	uid := ctx.Param("uid")
	if err := roleRepo.DeleteRole(uid); err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": "The Role has been deleted successfully"})
}

func GetRole(ctx *gin.Context) {
	uid := ctx.Param("uid")
	role, err := roleRepo.FindRoleByUid(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, role)
}
