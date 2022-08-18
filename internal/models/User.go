package models

import (
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/requests"
)

type User struct {
	ID         uint   `json:"id" gorm:"primary key"`
	Name       string `json:"name" gorm:"unique"`
	TgID       uint   `json:"tg_id" gorm:"unique"`
	Categories []Category
}

func (u *User) Create(db *database.SqlLite, req requests.User) error {
	u.Name = req.Name
	u.TgID = req.TgID
	result := db.DB.Select("Name", "TgID").Create(u)

	return result.Error
}

func (u *User) Get(db *database.SqlLite, id uint) error {
	result := db.DB.First(u, id)

	return result.Error
}

func (u *User) Update(db *database.SqlLite, id uint) error {
	name := u.Name
	err := u.Get(db, id)

	if err != nil {
		return err
	}

	result := db.DB.Model(u).Update("name", name)

	return result.Error
}

func (u *User) Delete(db *database.SqlLite, id uint) error {
	result := db.DB.Delete(u, id)

	return result.Error
}
