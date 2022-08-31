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

	if lastMessage.Message.Entities != nil {
		message.Command.String, message.Command.Valid = lastMessage.Message.Entities[0].Type, true
	}

	message.IsBot = lastMessage.Message.From.IsBot
	err := message.Create(router.DB)

	return err
}

func (router Router) saveHandledMsg(lastMessage telegram.MessageWrapper) error {
	botMsg, err := router.TgService.HandleMethods(lastMessage)

	if err != nil {
		return err
	}

	err = router.saveBotMsg(botMsg)

	return err
}

func (router Router) handleLastCommand(msg models.Message, lastMessage telegram.MessageWrapper) error {
	var err error

	switch msg.Text {
	case "/categoryCreate":
		var user models.User
		user.GetFromTg(router.DB, lastMessage.Message.From.Id)

		var category models.Category
		category.Name = lastMessage.Message.Text
		category.UserID = user.ID

		err = category.Create(router.DB)

		if err != nil {
			break
		}

		botMsg, err := router.TgService.SendCreated(lastMessage)

		if err != nil {
			break
		}

		err = router.saveBotMsg(botMsg)
		break
	}
	return err
}

func (router Router) saveBotMsg(botMsg telegram.BotMessage) error {
	var message models.Message
	message.Text = botMsg.Result.Text
	message.TgID = uint(botMsg.Result.MessageId)
	message.UserId = uint(botMsg.Result.Chat.Id)
	message.IsBot = botMsg.Result.From.IsBot

	err := message.Create(router.DB)

	return err
}

func (router *Router) AnswerToCallback(lastUpdate telegram.MessageWrapper) {
	if lastUpdate.CallbackQuery.Id != "" {
		router.TgService.AnswerCallbackQuery(lastUpdate.CallbackQuery.Id, "Just wait")
	}
}
