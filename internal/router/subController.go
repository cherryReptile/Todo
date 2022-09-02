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

func (router *Router) HandleMethods(message telegram.MessageWrapper) (telegram.BotMessage, error) {
	var botMsg telegram.BotMessage
	var err error

	switch message.Message.Text {
	case "/start":
		botMsg, err = router.TgService.SendHello(message)
		break
	case "/categoryCreate":
		botMsg, err = router.TgService.SendCreate(message)
		break
	case "/list":
		err = router.List(message, "Твои категории(нажми чтобы увидеть todo): 👇\n")
		break
	case "/categoryDelete":
		err = router.List(message, "Выберите какую категорию удалить 🗑\n")
	default:
		botMsg, err = router.TgService.SendDefault(message)
		break
	}
	return botMsg, err
}

func (router Router) saveHandledMsg(lastMessage telegram.MessageWrapper) error {
	botMsg, err := router.HandleMethods(lastMessage)

	if err != nil {
		return err
	}

	err = router.saveBotMsg(botMsg)

	return err
}

func (router Router) handleLastCommand(msg models.Message, lastMessage telegram.MessageWrapper) error {
	var err error

	switch {
	case lastMessage.Message.Text == "/start":
		router.TgService.SendHello(lastMessage)
		break
	case lastMessage.Message.Text == "/categoryCreate":
		router.TgService.SendCreate(lastMessage)
		break
	case msg.Text == "/categoryCreate" && uint(lastMessage.Message.MessageId)-msg.TgID == 2:
		err = router.Create(lastMessage)
		break
	case lastMessage.Message.Text == "/list":
		err = router.List(lastMessage, "Твои категории(нажми чтобы увидеть todo): 👇\n")
		break
	case msg.Text == "/list" && lastMessage.CallbackQuery.Id != "":
		err = router.Get(lastMessage)
		break
	case lastMessage.Message.Text == "/categoryDelete":
		err = router.List(lastMessage, "Выберите какую категорию удалить 🗑\n")
		break
	case msg.Text == "/categoryDelete" && lastMessage.CallbackQuery.Id != "":
		err = router.Delete(lastMessage)
	default:
		botMsg, err := router.TgService.SendDefault(lastMessage)

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
