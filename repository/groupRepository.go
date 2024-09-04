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
	return repo.db.Where("uid=?", group.Uid).Save(group).Error
}

func (repo *GroupRepository) DeleteGroup(uid string) error {
	return repo.db.Where("uid = ?", uid).Delete(&models.Group{}).Error
}

func (repo *GroupRepository) FindGroupByUid(uid string) (*models.Group, error) {
	var group models.Group
	if err := repo.db.Where("uid = ?", uid).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (repo *GroupRepository) FindAllGroups() ([]*models.Group, error) {
	var groups []*models.Group
	err := repo.db.Find(&groups).Error
	return groups, err
}
