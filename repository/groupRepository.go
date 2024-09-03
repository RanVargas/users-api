package repository

import (
	"gorm.io/gorm"
	"users-api/models"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (repo *GroupRepository) CreateGroup(group *models.Group) error {
	return repo.db.Create(group).Error
}

func (repo *GroupRepository) UpdateGroup(group *models.Group) error {
	return repo.db.Save(group).Error
}

func (repo *GroupRepository) DeleteGroup(uid string) error {
	return repo.db.Delete("uid = ?", uid).Error
}

func (repo *GroupRepository) FindGroupById(uid string) error {
	return repo.db.Where("uid = ?", uid).Error
}

func (repo *GroupRepository) FindAllGroups() ([]*models.Group, error) {
	var groups []*models.Group
	err := repo.db.Find(&groups).Error
	return groups, err
}
