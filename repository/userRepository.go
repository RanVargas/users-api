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

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {

	return repo.db.Create(&user).Error
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.db.Model(&models.User{}).
		Where("email = ?", user.Email).
		Update("password", user.Password).Error
}

func (repo *UserRepository) DeleteUser(uid string) error {
	return repo.db.Where("uid = ?", uid).Delete(&models.User{}).Error
}

func (repo *UserRepository) GetUser(uid string) (*models.User, error) {
	var user models.User
	err := repo.db.Preload("Role").
		Where("uid = ?", uid).First(&user).Error
	return &user, err
}

func (repo *UserRepository) GetUserById(id int) (*models.User, error) {
	var user models.User
	if err := repo.db.Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (repo *UserRepository) GetAllUsers() ([]*models.User, error) {
	var user []*models.User
	err := repo.db.Preload("Role").
		Find(&user).Error
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
		"groups":  true,
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
	if err := repo.db.
		Joins("INNER JOIN groups_users_map ON groups_users_map.groups_id = groups.id").
		Joins("INNER JOIN users ON users.id = groups_users_map.users_id").
		Where("users.uid = ?", uid).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (repo *UserRepository) GetAllUsersByRoleId(roleId string) ([]*models.User, error) {
	var users []*models.User
	if err := repo.db.
		Preload("Role").
		Joins("INNER JOIN roles ON roles.id = users.role_id").
		Where("roles.uid = ?", roleId).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *UserRepository) UpdateUserPassword(user *models.User) error {
	return repo.db.Model(user).Update("password", user.Password).Where("uid = ?", user.Uid).Error
}
