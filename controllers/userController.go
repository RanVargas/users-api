package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"users-api/database"
	"users-api/models"
	"users-api/repository"
)

var userRepo *repository.UserRepository

var userPasswordRepo *repository.UserPasswordRepository

func InitializeUsersRepo() {
	userRepo = repository.NewUserRepository(database.DB)
}

func InitializeUserPasswordRepo() {
	userPasswordRepo = repository.NewUserPasswordRepository(database.DB)
}

func CreateUser(ctx *gin.Context) {
	var body struct {
		Name      string   `json:"name" binding:"required"`
		Email     string   `json:"email" binding:"required"`
		Password  string   `json:"password" binding:"required"`
		Status    int16    `json:"status" binding:"required"`
		RoleUid   string   `json:"role_uid" binding:"required"`
		GroupsUid []string `json:"groups_uid" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	body.Password = string(hash)
	role, err := roleRepo.FindRoleByUid(body.RoleUid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong Role Uid provided"})
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
		Email:    body.Email,
		Password: body.Password,
		Name:     body.Name,
		Status:   body.Status,
		RoleID:   role.ID,
		Group:    groups,
	}

	if err := userRepo.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user, details: " + err.Error()})
		return
	}

	userPassword := models.UserPassword{
		UserID:       user.ID,
		UserPassword: user.Password,
	}
	if err := userPasswordRepo.CreateUserPassword(&userPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": body})
}

func GetAllUsers(ctx *gin.Context) {

	result, err := userRepo.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := userRepo.UpdateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("uid")
	if err := userRepo.DeleteUser(id); err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": "The User has been deleted successfully"})
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func FindUsersByQueryParams(c *gin.Context, searchTerm string, limitParam int, orderBy string) {
	//var users []models.User
	users, err := userRepo.FindUsersByQueryParameters(searchTerm, limitParam, orderBy)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUserAndRoleByUid(ctx *gin.Context) {
	uid := ctx.Param("uid")
	user, err := userRepo.GetUserAndRoleByUid(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func GetUsersByRole(ctx *gin.Context) {
	roleId := ctx.Param("uid")
	users, err := userRepo.GetAllUsersByRoleId(roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

func GetGroupsOfUser(ctx *gin.Context) {
	uid := ctx.Param("uid")
	groups, err := userRepo.GetAllGroupsOfUser(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": groups})
}
