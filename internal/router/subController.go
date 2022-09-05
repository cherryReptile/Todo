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
		err = message.Create(router.DB)
	case lastMessage.CallbackQuery.Id != "":
		var callback models.Callback
		bytes, err := json.Marshal(&lastMessage.CallbackQuery)

		if err != nil {
			break
		}

		callback.Json = string(bytes)
		callback.TgID = uint(lastMessage.CallbackQuery.Message.MessageId)
		callback.UserId = uint(lastMessage.CallbackQuery.Chat.Id)
		err = callback.Create(router.DB)
	}

	return err
}

func (router Router) handleLastCommand(lastCommand models.Message, modelFromCallback telegram.ModelFromCallback, lastUpdate telegram.MessageWrapper) error {
	var err error

	switch {
	case lastUpdate.Message.Text == "/start":
		router.TgService.SendHello(lastUpdate)
		break
	case lastUpdate.Message.Text == "/categoryCreate":
		router.TgService.SendCreate(lastUpdate)
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
	case lastCommand.Text == "/todo" && modelFromCallback.Model == "category" || modelFromCallback.Model == "todo":
		err = router.TodoController.Create(lastUpdate, modelFromCallback)
		break
	//case lastCommand.Text == "/todo" && modelFromCallback.Model == "todo":
	//	err = router.TodoController.Delete(lastUpdate, modelFromCallback)
	//case lastCommand.Text == "/todoCreate":
	//	err = router.CategoryController.List(lastUpdate, "Выберите в какой категории создать todo ✍️\n")
	//case lastCommand.Text == "/todoCreate" && modelFromCallback.Model == "":

	default:
		botMsg, err := router.TgService.SendDefault(lastUpdate)

		if err != nil {
			break
		}

		err = controllers.SaveBotMsg(botMsg, router.DB)
		break
	}
	return err
}
