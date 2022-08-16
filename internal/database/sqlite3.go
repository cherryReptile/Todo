package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type SqlLite struct {
	DB *gorm.DB
}

func Connect() SqlLite {
	db, err := gorm.Open(sqlite.Open("./tmp/db.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return SqlLite{DB: db}
}
