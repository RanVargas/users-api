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

var groupRepo *repository.GroupRepository

func InitializeGroupsRepo() {
	groupRepo = repository.NewGroupRepository(database.DB)
}

func CreateGroup(ctx *gin.Context) {
	var group *models.Group
	if err := ctx.ShouldBindJSON(&group); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := groupRepo.CreateGroup(group)
	if err != nil {
		log.Printf("Error saving to database the group, %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating group"})
		return
	}

	ctx.JSON(http.StatusCreated, group)
}

func GetAllGroups(ctx *gin.Context) {
	groups, err := groupRepo.FindAllGroups()
	if err != nil {
		log.Printf("Error retrieving from database the groups, %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving groups"})
		return
	}
	ctx.JSON(http.StatusOK, groups)
}

func UpdateGroup(ctx *gin.Context) {
	var group *models.Group
	if err := ctx.ShouldBindJSON(&group); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := groupRepo.UpdateGroup(group)
	if result != nil {
		log.Printf("Error while saving updated group %v", result)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error while updating group"})
		return
	}
	ctx.JSON(http.StatusOK, group)
}

func DeleteGroup(ctx *gin.Context) {
	id := ctx.Param("uid")
	if err := groupRepo.DeleteGroup(id); err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": "The Group has been deleted successfully"})
}

func GetGroup(ctx *gin.Context) {
	id := ctx.Param("uid")
	group, err := groupRepo.FindGroupByUid(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, group)
}
