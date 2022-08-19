package models

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/requests"
)

type User struct {
	ID   uint   `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	TgID uint   `json:"tg_id" db:"tg_id"`
}

func (u *User) Create(db *database.SqlLite, req *requests.User) error {
	u.Name = req.Name
	u.TgID = req.TgID
	result, err := db.DB.NamedExec("INSERT INTO users (name, tg_id) VALUES (:name, :tg_id)", u)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	u.ID = uint(id)
	fmt.Println(result)
	return err
}

func (u *User) Get(db *database.SqlLite, id uint) error {
	err := db.DB.Get(u, "SELECT * FROM users WHERE id=?", id)

	return err
}

func (u *User) Update(db *database.SqlLite, id uint) error {
	name := u.Name
	result, err := db.DB.Exec("UPDATE users SET name=? WHERE id=?", name, id)
	fmt.Println(result)

	if err != nil {
		return err
	}

	err = u.Get(db, id)

	return err
}

func (u *User) Delete(db *database.SqlLite, id uint) error {
	result, err := db.DB.Exec("DELETE FROM users WHERE id=?", id)
	fmt.Println(result)

	return err
}
