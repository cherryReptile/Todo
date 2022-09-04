package models

import (
	"database/sql"
	"github.com/cherryReptile/Todo/internal/database"
	"time"
)

type Callback struct {
	ID        uint         `db:"id"`
	Json      string       `db:"json"`
	TgID      uint         `db:"tg_id"`
	UserId    uint         `db:"user_id"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

func (c *Callback) Create(db *database.SqlLite) error {
	c.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	result, err := db.DB.NamedExec("INSERT INTO callbacks (json, tg_id, user_id, created_at) VALUES (:json, :tg_id, :user_id, :created_at)", c)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	c.Get(db, uint(id))

	return err
}

func (c *Callback) Get(db *database.SqlLite, id uint) error {
	err := db.DB.Get(c, "SELECT * FROM callbacks WHERE id=? LIMIT 1", id)

	return err
}

func (c *Callback) GetLast(db *database.SqlLite, userId uint) error {
	err := db.DB.Get(c, "SELECT * FROM callbacks WHERE user_id=? ORDER BY id DESC LIMIT 1", userId)

	return err
}
