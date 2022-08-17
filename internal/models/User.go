package models

type User struct {
	ID   uint   `json:"id" gorm:"primary key"`
	Name string `json:"name" gorm:"unique" validate:"required"`
	TgID uint   `json:"tg_id" gorm:"unique" validate:"required"`
}
