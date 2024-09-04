package models

import (
	"gorm.io/gorm"
)

type GroupUserMap struct {
	gorm.Model
	//Id      uint  `gorm:"primaryKey;autoIncrement"`
	UserID  uint  `json:"user_id"`
	User    User  `gorm:"foreignKey:UserID"`
	GroupID uint  `json:"group_id"`
	Group   Group `gorm:"foreignKey:GroupID"`
}

func (GroupUserMap) TableName() string {
	return "groups_users_map"
}
