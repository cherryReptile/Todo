package models

import (
	"database/sql"
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/requests"
	"time"
)

type Todo struct {
	ID         uint         `json:"id" db:"id"`
	Name       string       `json:"name" db:"name"`
	CategoryID uint         `json:"category_id" db:"category_id"`
	CreatedAt  sql.NullTime `json:"created_at" db:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at" db:"updated_at"`
}

func (t *Todo) ItsRequest() requests.Todo {
	return requests.Todo{}
}

func (t *Todo) Create(db *database.SqlLite, req *requests.Todo) error {
	t.Name = req.Name
	t.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	result, err := db.DB.NamedExec("INSERT INTO todos (name, category_id, created_at) VALUES (:name, :category_id, :created_at)", t)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	t.ID = uint(id)
	fmt.Println(result)
	return err
}

func (t *Todo) Get(db *database.SqlLite, id uint) error {
	err := db.DB.Get(t, "SELECT * FROM todos WHERE id=?", id)

	return err
}

func (t *Todo) Update(db *database.SqlLite, id uint) error {
	name := t.Name
	result, err := db.DB.Exec("UPDATE todos SET name=? WHERE id=?", name, id)
	fmt.Println(result)

	if err != nil {
		return err
	}

	err = t.Get(db, id)

	return err
}

func (t *Todo) Delete(db *database.SqlLite, id uint) error {
	result, err := db.DB.Exec("DELETE FROM todos WHERE id=?", id)
	fmt.Println(result)

	return err
}
