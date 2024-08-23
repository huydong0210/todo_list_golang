package model

import "gorm.io/gorm"

type TodoItem struct {
	gorm.Model
	Name  string
	State string
}

func CreateTodoItem(db *gorm.DB, item *TodoItem) error {
	result := db.Create(item)
	return result.Error
}
func DeleteTodoItem(db *gorm.DB, id int) error {
	result := db.Delete(&TodoItem{}, id)
	return result.Error
}
func UpdateTodoItem(db *gorm.DB, id int, item *TodoItem) error {
	result := db.Model(&TodoItem{}).Where("id = ?", id).Updates(item)
	return result.Error
}
