package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"users-api/database"
	"users-api/models"
	"users-api/repository"
)

var groupRepo *repository.GroupRepository

func InitializeGroupsRepo() {
	groupRepo = repository.NewGroupRepository(database.DB)
}

func CreateGroup(ctx *gin.Context) {
	var group models.Group
	if err := ctx.ShouldBindJSON(&group); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := database.DB.Create(&group)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": group})
}

func GetAllGroups(ctx *gin.Context) {
	var groups []models.Group
	result := database.DB.Find(&groups)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": groups})
}

func UpdateGroup(ctx *gin.Context) {
	var group models.Group
	if err := ctx.ShouldBindJSON(&group); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Save(&group)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": group})
}

func DeleteGroup(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := database.DB.Delete(&models.Group{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": "The Group has been deleted successfully"})
}

func GetGroup(ctx *gin.Context) {
	id := ctx.Param("id")
	var group models.Group
	if err := database.DB.First(&group, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": group})
}
