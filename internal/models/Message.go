package models

import (
	"database/sql"
	"github.com/cherryReptile/Todo/internal/database"
	"time"
)

type Message struct {
	ID        uint           `json:"id" db:"id"`
	Text      string         `json:"text" db:"text"`
	TgID      uint           `json:"tg_id" db:"tg_id"`
	UserId    uint           `json:"user_id" db:"user_id"`
	IsBot     bool           `json:"is_bot" db:"is_bot"`
	Command   sql.NullString `json:"command" db:"command"`
	CreatedAt sql.NullTime   `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at" db:"updated_at"`
}

func (m *Message) Create(db *database.SqlLite) error {
	m.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	result, err := db.DB.NamedExec("INSERT INTO messages (text, tg_id, user_id, is_bot, command, created_at) VALUES (:text, :tg_id, :user_id, :is_bot,:command, :created_at)", m)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	m.Get(db, uint(id))

	return err
}

func (m *Message) Get(db *database.SqlLite, id uint) error {
	err := db.DB.Get(m, "SELECT * FROM messages WHERE id=? LIMIT 1", id)

	return err
}

func (m *Message) GetFromTg(db *database.SqlLite, tgId uint) error {
	err := db.DB.Get(m, "SELECT * FROM messages WHERE tg_id=? LIMIT 1", tgId)

	return err
}

func (m *Message) GetLastCommand(db *database.SqlLite, userId uint) error {
	m.UserId = userId
	rows, err := db.DB.NamedQuery("SELECT * FROM messages WHERE user_id=:user_id AND command='bot_command' ORDER BY id DESC LIMIT 1", m)
	for rows.Next() {
		err = rows.StructScan(m)
		if err != nil {
			return err
		}
	}
	return err
}
