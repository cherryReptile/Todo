package controllers

import (
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/telegram"
)

type DbController struct {
	DB *database.SqlLite
}

type TgController struct {
	TgService *telegram.Service
}
