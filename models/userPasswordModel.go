package models

import "gorm.io/gorm"

type UserPassword struct {
	gorm.Model `json:"-"`
	//Id           uint   `gorm:"primaryKey;autoIncrement" json:"-"`
	UserID       uint   `json:"-"`
	User         User   `gorm:"foreignKey:UserID" json:"-"`
	UserPassword string `gorm:"not null" json:"-"`
}
