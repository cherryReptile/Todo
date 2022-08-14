package models

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
)

type User struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (u *User) Create(db *database.SqlLite) error {
	result, err := db.DB.NamedExec("INSERT INTO users (name) VALUES (:name)", u)

	if err != nil {
		return err
	}

	u.ID, err = result.LastInsertId()
	fmt.Println(result)
	return err
}

func (u *User) Get(db *database.SqlLite, id int) error {
	rows, err := db.DB.Queryx(`SELECT * FROM users WHERE id=?`, id)

	for rows.Next() {
		err = rows.StructScan(u)
	}

	fmt.Println(u)
	return err
}

func (u *User) Update() {
	//
}
