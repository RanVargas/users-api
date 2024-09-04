package repository

import (
	"gorm.io/gorm"
	"users-api/models"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (repo *RoleRepository) FindRoleByUid(uid string) (*models.Role, error) {
	var role models.Role
	err := repo.db.Where("uid = ?", uid).First(&role).Error
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
	err := repo.db.Where("uid = ?", role.Uid).Save(&role).Error
	return err
}

func (repo *RoleRepository) DeleteRole(uid string) error {
	var role models.Role
	err := repo.db.Where("uid = ?", uid).Delete(&role).Error
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
