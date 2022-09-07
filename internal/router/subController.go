package router

import (
	"encoding/json"
	"github.com/cherryReptile/Todo/internal/controllers"
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

	if len(updates.Result) == 1 {
		return updates.Result[0], err
	}

	lastMessage := updates.Result[len(updates.Result)-1]
	return lastMessage, nil
}

func (router Router) saveIncomingMsg(lastMessage telegram.MessageWrapper) error {
	var err error

	switch {
	case lastMessage.Message.Text != "":
		var message models.Message

		message.Text = lastMessage.Message.Text
		message.TgID = uint(lastMessage.Message.MessageId)
		message.UserId = lastMessage.Message.From.Id

		if lastMessage.Message.Entities != nil {
			message.Command.String, message.Command.Valid = lastMessage.Message.Entities[0].Type, true
		}

		message.IsBot = lastMessage.Message.From.IsBot
		message.IsCallback = false
		err = message.Create(router.DB)
	case lastMessage.CallbackQuery.Id != "":
		var message models.Message
		bytes, err := json.Marshal(&lastMessage.CallbackQuery)

		if err != nil {
			break
		}

		message.Text = string(bytes)
		message.TgID = uint(lastMessage.CallbackQuery.Message.MessageId)
		message.UserId = uint(lastMessage.CallbackQuery.Chat.Id)
		message.IsBot = false
		message.IsCallback = true
		err = message.Create(router.DB)
	}

	return err
}

func (router Router) handleLastCommand(lastCommand models.Message, modelFromCallback telegram.ModelFromCallback, lastUpdate telegram.MessageWrapper) error {
	var err error
	var botMsg telegram.BotMessage

	switch {
	case lastUpdate.Message.Text == "/start":
		botMsg, err = router.TgService.SendHello(lastUpdate)
		break
	case lastUpdate.Message.Text == "/categoryCreate":
		botMsg, err = router.TgService.SendCreate(lastUpdate)
		break
	case lastCommand.Text == "/categoryCreate" && uint(lastUpdate.Message.MessageId)-lastCommand.TgID == 2:
		err = router.CategoryController.Create(lastUpdate)
		break
	case lastUpdate.Message.Text == "/list":
		err = router.CategoryController.List(lastUpdate, "Твои категории(нажми чтобы увидеть todo): 👇\n")
		break
	case lastCommand.Text == "/list" && modelFromCallback.Model == "category":
		err = router.CategoryController.Get(lastUpdate, modelFromCallback)
		break
	case lastCommand.Text == "/list" && modelFromCallback.Model == "todo":
		err = router.TodoController.Delete(lastUpdate, modelFromCallback)
	case lastUpdate.Message.Text == "/categoryDelete":
		err = router.CategoryController.List(lastUpdate, "Выберите какую категорию удалить 🗑:\n")
		break
	case lastCommand.Text == "/categoryDelete" && lastUpdate.CallbackQuery.Id != "":
		err = router.CategoryController.Delete(lastUpdate)
		break
	case lastUpdate.Message.Text == "/todo":
		err = router.CategoryController.List(lastUpdate, "Выберите в какой категории создать todo ✍️\n")
		break
	case lastCommand.Text == "/todo" && lastUpdate.CallbackQuery.Id != "":
		err = router.TodoController.Create(lastUpdate, modelFromCallback)
		break
	case lastCommand.Text == "/todo" && lastUpdate.CallbackQuery.Id == "":
		err = router.TodoController.DefaultCreate(lastUpdate, modelFromCallback)
	default:
		botMsg, err = router.TgService.SendDefault(lastUpdate)
		break
	}

	if err != nil {
		return err
	}

	if botMsg.Result.Text != "" {
		err = controllers.SaveBotMsg(botMsg, router.DB)
	}

	if err != nil {
		return err
	}

	return nil
}
