package models

import (
	"gorm.io/gorm"
)

type GroupUserMap struct {
	gorm.DB
	Id      uint  `gorm:"primaryKey;autoIncrement"`
	UserID  int   `gorm:"column:User_id" json:"user_id"`
	GroupID int   `gorm:"column:Group_id" json:"group_id"`
	User    User  `gorm:"foreignKey:UserID"`
	Group   Group `gorm:"foreignKey:GroupID"`
}

func (GroupUserMap) TableName() string {
	return "groups_users_map"
}
