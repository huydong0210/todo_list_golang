package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
	Email    string
	RoleId   int
}

func FindByUsername(db *gorm.DB, username string) error {
	var user User
	result := db.Where("username = ?", username).First(&user)
	return result.Error
}
func DeleteByUsername(db *gorm.DB, username string) error {
	result := db.Where("username = ?", username).Delete(&User{})
	return result.Error
}
