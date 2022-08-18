package models

import "github.com/cherryReptile/Todo/internal/requests"

type User struct {
	ID   uint   `json:"id" gorm:"primary key"`
	Name string `json:"name" gorm:"unique"`
	TgID uint   `json:"tg_id" gorm:"unique"`
}

func NewUser(reqU requests.User) User {
	return User{
		ID:   0,
		Name: reqU.Name,
		TgID: reqU.TgID,
	}
}
