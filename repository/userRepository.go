package repository

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"users-api/models"
)

type UserRepository struct {
	db *gorm.DB
}

/*
/roles/uid/users : all users with a specific role
--------------/users : all users
-------------/users?search=searchterm_for_name&limit=max_number_of_results&orderby=column_order_by
------------/users/uid : user details with specific uid, including role
-----------/users/uid/groups : all groups of user #uid
*/

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.db.Save(user).Error
}

func (repo *UserRepository) DeleteUser(id string) error {
	return repo.db.Where("id = ?", id).Error
}

func (repo *UserRepository) GetUser(id string) (*models.User, error) {
	var user models.User
	err := repo.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (repo *UserRepository) GetAllUsers() ([]*models.User, error) {
	var user []*models.User
	err := repo.db.Find(&user).Error
	return user, err
}

func (repo *UserRepository) FindUsersByQueryParameters(search string, limit int, orderBy string) ([]models.User, error) {
	var users []models.User

	validColumns := map[string]bool{
		"id":      true,
		"name":    true,
		"email":   true,
		"uid":     true,
		"role_id": true,
	}

	if !validColumns[orderBy] {
		return nil, fmt.Errorf("invalid order by %s", orderBy)
	}

	query := repo.db.Model(&models.User{})
	if search != "" {
		query = repo.db.Where("name LIKE ?", "%"+search+"%").
			Or("name LIKE ?", "%"+strings.ToLower(search)+"%").
			Or("name LIKE ?", "%"+strings.ToUpper(search)+"%")
	}

	query = query.Order(orderBy).Limit(limit)
	err := query.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	return users, nil
}

func (repo *UserRepository) GetUserAndRoleByUid(uid string) (*models.User, error) {
	var user models.User
	if err := repo.db.
		Preload("Role").
		Joins("INNER JOIN roles ON roles.id = users.role_id").
		Where("users.uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetAllGroupsOfUser(uid string) ([]*models.Group, error) {
	var groups []*models.Group
	if err := repo.db.Joins("INNER JOIN groups_users_map ON groups_users_map.users_id = users.id").
		Joins("INNER JOIN groups ON groups.id = groups_users_map.group_id").
		//Where("groups_users_map.group_id = groups.id").
		Where("groups_users_map.users_id = ?", uid).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (repo *UserRepository) GetAllUsersByRoleId(roleId string) ([]*models.User, error) {
	var users []*models.User
	if err := repo.db.Where("role_id = ?", roleId).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetAllUsersByRoleId GetAllGroupsOfUser GetUserAndRoleByUid FindUsersByQueryParameters
