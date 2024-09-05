package models

import (
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model `json:"-"`
	Name       string          `gorm:"type:varchar(250);not null" json:"name"`
	Uid        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"uid"`
	Rights     json.RawMessage `gorm:"type:jsonb" json:"rights"`
}

func (Role) TableName() string {
	return "roles"
}
