package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
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
		log.Printf("Error creating user in database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user due to internal error"})
		return
	}

	userPassword := models.UserPassword{
		UserID:       user.ID,
		UserPassword: user.Password,
	}
	if err := userPasswordRepo.CreateUserPassword(&userPassword); err != nil {
		log.Printf("Error creating user password in database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user password due to internal error"})
		return
	}

	ctx.JSON(http.StatusCreated, body)
}

func GetAllUsers(ctx *gin.Context) {

	result, err := userRepo.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		log.Printf("Could not bind body to model: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not bind body to model, bad request data"})
		return
	}

	if err := userRepo.UpdateUser(&user); err != nil {
		log.Printf("Error updating user in database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user due to internal error"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("uid")
	if err := userRepo.DeleteUser(id); err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			log.Printf("Error deleting user in database: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not delete user due to internal error"})
		}
		log.Printf("Error deleting user in database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user due to internal error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": "The User has been deleted successfully"})
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
	ctx.JSON(http.StatusOK, user)
}

func FindUsersByQueryParams(c *gin.Context, searchTerm string, limitParam int, orderBy string) {
	users, err := userRepo.FindUsersByQueryParameters(searchTerm, limitParam, orderBy)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error finding users in database: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		log.Printf("Error finding users in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserAndRoleByUid(ctx *gin.Context) {
	uid := ctx.Param("uid")
	user, err := userRepo.GetUserAndRoleByUid(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error finding user in database: %v", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func GetUsersByRole(ctx *gin.Context) {
	roleId := ctx.Param("uid")
	users, err := userRepo.GetAllUsersByRoleId(roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error finding users in database: %v", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find users due to internal error"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func GetGroupsOfUser(ctx *gin.Context) {
	uid := ctx.Param("uid")
	groups, err := userRepo.GetAllGroupsOfUser(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error finding users in database: %v", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, groups)
}

func UpdateUserPassword(ctx *gin.Context) {
	uid := ctx.Param("uid")
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Uid      string `json:"uid" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("Encounterd Error while binding body to model: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request body"})
		return
	}
	user, er := userRepo.GetUser(uid)
	if er != nil {
		if errors.Is(er, gorm.ErrRecordNotFound) {
			log.Printf("Encounterd Error getting user: %v", er)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User to Update Not Found"})
			return
		}
		log.Printf("Encounterd Error while retrieving user %v", er)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	userPassword, err := userPasswordRepo.GetUserPassword(user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Encounterd Error while getting user-password: %v", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User-password to Update Not Found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error " + err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		log.Printf("Encounterd Error while encrypting user password: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error encrypting password"})
		return
	}
	body.Password = string(hash)
	userPassword.UserPassword = body.Password
	user.Password = body.Password
	if e := userPasswordRepo.UpdateUserPassword(userPassword); e != nil {
		log.Printf("Encounterd Error while updating user password: %v", e)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving password to database"})
		return
	}
	if e := userRepo.UpdateUserPassword(user); e != nil {
		log.Printf("Encounterd Error while updating user password: %v", e)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving password to database"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})

}
