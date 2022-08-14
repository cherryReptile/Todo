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

//func (u *User) Get(db *database.SqlLite) (sql.Result, error) {
//	result, err := db.DB.NamedExec(`SELECT * FROM users WHERE id=:id`, u)
//
//	fmt.Println(result)
//	return result, err
//}
