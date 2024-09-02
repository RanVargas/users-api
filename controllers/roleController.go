package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"users-api/database"
	"users-api/models"
)

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

	ctx.JSON(http.StatusCreated, gin.H{"data": role})
}

func GetAllRoles(ctx *gin.Context) {
	var roles []models.Role
	result := database.DB.Find(&roles)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": roles})
}

func UpdateRole(ctx *gin.Context) {
	var role models.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Save(&role)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": role})
}

func DeleteRole(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := database.DB.Delete(&models.Role{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": "The Role has been deleted successfully"})
}

func GetRole(ctx *gin.Context) {
	id := ctx.Param("id")
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": role})
}
