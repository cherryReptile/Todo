package models

import (
	"database/sql"
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"time"
)

type Category struct {
	ID        uint         `json:"id" db:"id"`
	Name      string       `json:"name" db:"name"`
	UserID    uint         `json:"user_id" db:"user_id"`
	CreatedAt sql.NullTime `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"`
}

func (c *Category) Create(db *database.SqlLite) error {
	c.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	result, err := db.DB.NamedExec("INSERT INTO categories (name, user_id, created_at) VALUES (:name, :user_id, :created_at)", c)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	c.Get(db, uint(id))

	return err
}

func (c *Category) Get(db *database.SqlLite, id uint) error {
	err := db.DB.Get(c, "SELECT * FROM categories WHERE id=?", id)

	return err
}

func (c *Category) Update(db *database.SqlLite, id uint) error {
	name := c.Name
	result, err := db.DB.Exec("UPDATE categories SET name=? WHERE id=?", name, id)
	fmt.Println(result)

	if err != nil {
		return err
	}

	err = c.Get(db, id)

	return err
}

func (c *Category) Delete(db *database.SqlLite, id uint) error {
	result, err := db.DB.Exec("DELETE FROM categories WHERE id=?", id)
	fmt.Println(result)

	return err
}
