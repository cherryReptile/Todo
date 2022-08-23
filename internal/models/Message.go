package models

import (
	"database/sql"
	"github.com/cherryReptile/Todo/internal/database"
	"time"
)

type Message struct {
	ID        uint           `json:"id" db:"id"`
	Text      string         `json:"text" db:"text"`
	MessageID uint           `json:"Message_id" db:"message_id"`
	UserId    uint           `json:"user_id" db:"user_id"`
	Command   sql.NullString `json:"command" db:"command"`
	CreatedAt sql.NullTime   `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at" db:"updated_at"`
}

func (m *Message) Create(db *database.SqlLite) error {
	m.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	result, err := db.DB.NamedExec("INSERT INTO messages (text, message_id, user_id, command, created_at) VALUES (:name, :user_id, :created_at)", m)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	m.Get(db, uint(id))

	return err
}

func (m *Message) Get(db *database.SqlLite, id uint) error {
	err := db.DB.Get(m, "SELECT * FROM messages WHERE id=?", id)

	return err
}
