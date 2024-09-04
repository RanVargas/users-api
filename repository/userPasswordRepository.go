package repository

import (
	"gorm.io/gorm"
	"log"
	"users-api/models"
)

type UserPasswordRepository struct {
	db *gorm.DB
}

func NewUserPasswordRepository(db *gorm.DB) *UserPasswordRepository {
	return &UserPasswordRepository{db: db}
}

func (repo *UserPasswordRepository) CreateUserPassword(userPassword *models.UserPassword) error {
	log.Printf("This is the value of the user id from userpassword %v", userPassword.UserID)
	return repo.db.Preload("Users").Create(userPassword).Error
}

func (repo *UserPasswordRepository) GetUserPassword(userID uint) (*models.UserPassword, error) {
	var userPassword models.UserPassword //repo.db.Where("uid = ?", uid).First(&user).Error
	if err := repo.db.Where("user_id", userID).First(&userPassword).Error; err != nil {
		return nil, err
	}
	return &userPassword, nil
}
