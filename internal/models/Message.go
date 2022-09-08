package models

import (
	"database/sql"
	"github.com/cherryReptile/Todo/internal/database"
	"time"
)

type Message struct {
	ID         uint           `json:"id" db:"id"`
	Text       string         `json:"text" db:"text"`
	TgID       uint           `json:"tg_id" db:"tg_id"`
	UserId     uint           `json:"user_id" db:"user_id"`
	IsBot      bool           `json:"is_bot" db:"is_bot"`
	IsCallback bool           `json:"is_callback" db:"is_callback"`
	Command    sql.NullString `json:"command" db:"command"`
	CreatedAt  sql.NullTime   `json:"created_at" db:"created_at"`
	UpdatedAt  sql.NullTime   `json:"updated_at" db:"updated_at"`
}

func (m *Message) Create(db *database.SqlLite) error {
	m.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	result, err := db.DB.NamedExec("INSERT INTO messages (text, tg_id, user_id, is_bot, is_callback, command, created_at) VALUES (:text, :tg_id, :user_id, :is_bot, :is_callback,:command, :created_at)", m)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	m.Get(db, uint(id))

	return err
}

func (m *Message) Get(db *database.SqlLite, id uint) error {
	err := db.DB.Get(m, "SELECT * FROM messages WHERE id=$1 LIMIT 1", id)

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

func (m *Message) GetLastTwo(db *database.SqlLite, userId uint) ([]Message, error) {
	var msgs []Message
	err := db.DB.Select(&msgs, "SELECT * FROM messages WHERE user_id=:user_id AND is_bot=false ORDER BY id DESC LIMIT 2", userId)

	return msgs, err
}

func (m *Message) GetLast(db *database.SqlLite, userId uint) error {
	err := db.DB.Get(m, "SELECT * FROM messages WHERE user_id=$1 ORDER BY id DESC LIMIT 1", userId)

	return err
}

func (m *Message) GetLastBot(db *database.SqlLite, userId uint) error {
	err := db.DB.Get(m, "SELECT * FROM messages WHERE user_id=$1 AND is_bot=true ORDER BY id DESC LIMIT 1", userId)

	return err
}

func (m *Message) GetLastCallback(db *database.SqlLite, userId uint) error {
	err := db.DB.Get(m, "SELECT * FROM messages WHERE user_id=$1 AND is_callback=true ORDER BY id DESC LIMIT 1", userId)

	return err
}
