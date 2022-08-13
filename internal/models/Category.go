package models

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"time"
)

type Category struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	UserID    int       `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"Updated_at" db:"updated_at"`
}

func (c *Category) Create(db *database.SqlLite) error {
	c.CreatedAt = time.Now()

	result, err := db.DB.NamedExec("INSERT INTO categories(name, user_id, created_at) VALUES(:name, :user_id)", c)

	if err != nil {
		return err
	}

	c.ID, err = result.LastInsertId()
	fmt.Println(result)
	return err
}
