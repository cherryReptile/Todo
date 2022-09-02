package controllers

import (
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/telegram"
)

//Methods for controllers and router controller

func SaveBotMsg(botMsg telegram.BotMessage, db *database.SqlLite) error {
	var message models.Message
	message.Text = botMsg.Result.Text
	message.TgID = uint(botMsg.Result.MessageId)
	message.UserId = uint(botMsg.Result.Chat.Id)
	message.IsBot = botMsg.Result.From.IsBot

	err := message.Create(db)

	return err
}

func AnswerToCallback(lastUpdate telegram.MessageWrapper, TgService *telegram.Service) {
	if lastUpdate.CallbackQuery.Id != "" {
		TgService.AnswerCallbackQuery(lastUpdate.CallbackQuery.Id, "Just wait")
	}
}
