package models

import (
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model `json:"-"`
	Id         uint      `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name       string    `json:"name"`
	uid        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
}

func (Group) TableName() string {
	return "groups"
}
