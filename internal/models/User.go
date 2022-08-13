package models

import (
	"fmt"
	"github.com/pavel-one/GoStarter/internal/database"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (u *User) Create(db *database.SqlLite) error {
	result, err := db.DB.NamedExec("INSERT INTO users (name) VALUES (:name)", u)

	fmt.Println(result)
	return err
}
