package models

import (
	"gorm.io/gorm"
)

type GroupUserMap struct {
	gorm.Model
	UserID  uint  `json:"users_id"`
	User    User  `gorm:"foreignKey:UserID"`
	GroupID uint  `json:"groups_id"`
	Group   Group `gorm:"foreignKey:GroupID"`
}

func (GroupUserMap) TableName() string {
	return "groups_users_map"
}
