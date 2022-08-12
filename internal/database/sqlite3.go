package database

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type SqlLite struct {
	DB *sqlx.DB
}

func Connect() SqlLite {
	db, err := sqlx.Open("sqlite3", "./tmp/db.db")
	if err != nil {
		log.Fatal(err)
	}

	return SqlLite{DB: db}
}
