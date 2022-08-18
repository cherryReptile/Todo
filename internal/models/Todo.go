package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Name       string `json:"name"`
	CategoryID uint
}
