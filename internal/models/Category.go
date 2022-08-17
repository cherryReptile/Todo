package models

import (
	"time"
)

type Category struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	UserID    int       `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"Updated_at" db:"updated_at"`
}
