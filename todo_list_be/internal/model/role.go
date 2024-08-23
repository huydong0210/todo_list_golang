package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null; unique"`
}
