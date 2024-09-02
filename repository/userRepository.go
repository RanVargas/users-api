package repository

import (
	"gorm.io/gorm"
	"users-api/models"
)

type UserRepository struct {
	db *gorm.DB
}

func (repo *UserRepository) NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.db.Save(user).Error
}

func (repo *UserRepository) DeleteUser(id string) error {
	return repo.db.Where("id = ", id).Error
}

func (repo *UserRepository) GetUser(id string) (*models.User, error) {
	var user models.User
	err := repo.db.Where("id = ", id).First(&user).Error
	return &user, err
}

func (repo *UserRepository) GetAllUsers() (*models.User, error) {
	var user models.User
	err := repo.db.Find(&user).Error
	return &user, err
}
