package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name   string `json:"name"`
	UserID uint   `json:"user_id"`
	Todos  []Todo
}
