package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model `json:"-"`
	Id         uint              `gorm:"primaryKey;autoIncrement"`
	Name       string            `gorm:"type:varchar(250);not null" json:"name"`
	Uid        uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4()"`
	Rights     pgtype.JSONBCodec `gorm:"type:jsonb" json:"rights"`
}

func (Role) TableName() string {
	return "roles"
}
