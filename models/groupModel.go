package models

import (
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model `json:"-"`
	Name       string    `json:"name"`
	Uid        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Users      []User    `gorm:"many2many:groups_users_map;foreignKey:ID;joinForeignKey:group_id;References:ID;joinReferences:user_id"`
}

func (Group) TableName() string {
	return "groups"
}
