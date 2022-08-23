package router

import (
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/telegram"
)

func (router *Router) getUpdates() (telegram.Updates, error) {
	updates, err := router.TgService.GetUpdates()
	return updates, err
}

func (router *Router) getLastMsg() (telegram.MessageWrapper, error) {
	updates, err := router.getUpdates()

	if err != nil {
		return telegram.MessageWrapper{}, err
	}

	lastMessage := updates.Result[len(updates.Result)-1]
	return lastMessage, nil
}

func (router Router) saveIncomingMsg(lastMessage telegram.MessageWrapper) error {
	var message models.Message

	message.Text = lastMessage.Message.Text
	message.TgID = uint(lastMessage.Message.MessageId)
	message.UserId = lastMessage.Message.From.Id
	message.Command.String, message.Command.Valid = lastMessage.Message.Entities[0].Type, true
	message.IsBot = lastMessage.Message.From.IsBot
	err := message.Create(router.DB)

	return err
}

func (router Router) saveOutgoingMsg(lastMessage telegram.MessageWrapper) error {
	botMsg, err := router.TgService.HandleMethods(lastMessage)

	if err != nil {
		return err
	}

	var message models.Message
	message.Text = botMsg.Result.Text
	message.TgID = uint(botMsg.Result.MessageId)
	message.UserId = lastMessage.Message.From.Id
	message.IsBot = botMsg.Result.From.IsBot
	err = message.Create(router.DB)

	return err
}
