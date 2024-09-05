package models

import (
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Name       string    `gorm:"type:varchar(250)" json:"name"`
	Password   string    `gorm:"type:varchar(255)" json:"password"`
	Uid        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email      string    `gorm:"type:varchar(100);unique" json:"email"`
	Status     int16     `gorm:"default:0" json:"status"`
	RoleID     uint      `json:"role_id"`
	Role       Role      `gorm:"foreignKey:RoleID" json:"role"`
	Group      []Group   `gorm:"many2many:groups_users_map;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:group_id" json:"groups"`
}

func (User) TableName() string {
	return "users"
}
