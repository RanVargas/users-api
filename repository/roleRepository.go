package repository

import (
	"gorm.io/gorm"
	"users-api/models"
)

type RoleRepository struct {
	db *gorm.DB
}

func (repo *RoleRepository) NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (repo *RoleRepository) FindRoleById(id uint) (*models.Role, error) {
	var role models.Role
	err := repo.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (repo *RoleRepository) CreateRole(role models.Role) (*models.Role, error) {
	err := repo.db.Create(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (repo *RoleRepository) UpdateRole(role models.Role) error {
	err := repo.db.Save(&role).Error
	return err
}

func (repo *RoleRepository) DeleteRole(id uint) error {
	var role models.Role
	err := repo.db.Where("id = ?", id).Delete(&role).Error
	return err
}

func (repo *RoleRepository) FindAllRole() (*[]models.Role, error) {
	var roles []models.Role
	err := repo.db.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return &roles, nil
}
