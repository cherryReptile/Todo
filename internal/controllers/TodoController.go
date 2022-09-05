package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/telegram"
)

type TodoController struct {
	DB        *database.SqlLite
	TgService *telegram.Service
}

func NewTodoController(db *database.SqlLite, service *telegram.Service) *TodoController {
	return &TodoController{
		DB:        db,
		TgService: service,
	}
}

func (t *TodoController) Create(lastMessage telegram.MessageWrapper, modelFromCallback telegram.ModelFromCallback) error {
	go AnswerToCallback(lastMessage, t.TgService)

	var category models.Category
	category.Get(t.DB, modelFromCallback.Id)

	var todo models.Todo
	var err error

	todos, err := todo.GetAllFromCategoryId(t.DB, category.ID)

	if err != nil {
		return err
	}

	switch {
	case modelFromCallback.Model == "category" && lastMessage.CallbackQuery.Id != "":
		text := fmt.Sprintf("Todo %v категории(нажми чтобы удалить)\n", category.Name)

		if todos == nil {
			text := "У этой категории нет todo(вводите названия для создания)"
			_, err = t.TgService.EditMessageText(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, text)
			return err
		}

		_, err = t.TgService.EditMessageText(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, text)

		if err != nil {
			break
		}

		_, err = t.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, todos)
	case modelFromCallback.Model == "todo" && lastMessage.CallbackQuery.Id != "":
		err = t.Delete(lastMessage, modelFromCallback)
	//case modelFromCallback.Model == "category" && lastMessage.CallbackQuery.Id == "":
	//	todo.Name = lastMessage.Message.Text
	//	todo.CategoryID = category.ID
	//	err := todo.Create(t.DB)
	//
	//	if err != nil {
	//		break
	//	}
	//
	//	todos, err := todo.GetAllFromCategoryId(t.DB, category.ID)
	//
	//	if err != nil {
	//		break
	//	}
	//
	//	if len(todos) == 1 {
	//		text := fmt.Sprintf("Todo %v категории(нажми чтобы удалить)\n", category.Name)
	//		botMsg, err := t.TgService.SendInlineKeyboard(text, lastMessage.Message.From.Id, todos)
	//
	//		if err != nil {
	//			return err
	//		}
	//
	//		err = SaveBotMsg(botMsg, t.DB)
	//		fmt.Println(botMsg)
	//		return err
	//	}
	//
	//	var lastMsg models.Message
	//	lastMsg.GetLastBot(t.DB, lastMessage.Message.From.Id)
	//
	//	_, err = t.TgService.EditMessageReplyMarkup(lastMessage.Message.From.Id, int(lastMsg.TgID), todos)
	default:
		botMsg, err := t.TgService.SendDefault(lastMessage)

		if err != nil {
			break
		}

		err = SaveBotMsg(botMsg, t.DB)
		break
	}

	if err != nil {
		return err
	}

	return nil
}

func (t *TodoController) DefaultCreate(lastMessage telegram.MessageWrapper, modelFromCallback telegram.ModelFromCallback) error {
	var lastCallback models.Callback
	lastCallback.GetLast(t.DB, lastMessage.Message.From.Id)
	var callbackQuery telegram.CallbackQuery
	err := json.Unmarshal([]byte(lastCallback.Json), &callbackQuery)

	if err != nil {
		return err
	}

	if modelFromCallback.Model != "category" {
		botMsg, err := t.TgService.SendDefault(lastMessage)

		if err != nil {
			return err
		}

		err = SaveBotMsg(botMsg, t.DB)
		return err
	}

	var todo models.Todo
	todo.Name = lastMessage.Message.Text
	var category models.Category
	category.Get(t.DB, modelFromCallback.Id)

	if category.ID == 0 {
		err := errors.New("[ERROR]unknown category")
		return err
	}

	todo.CategoryID = category.ID
	err = todo.Create(t.DB)

	if err != nil {
		return err
	}

	todos, err := todo.GetAllFromCategoryId(t.DB, category.ID)

	if err != nil {
		return err
	}

	if len(todos) == 1 {
		text := fmt.Sprintf("Todo %v категории(нажми чтобы удалить)\n", category.Name)
		botMsg, err := t.TgService.SendInlineKeyboard(text, lastMessage.Message.From.Id, todos)

		if err != nil {
			return err
		}

		err = SaveBotMsg(botMsg, t.DB)
		fmt.Println(botMsg)
		return err
	}

	var lastMsg models.Message
	lastMsg.GetLastBot(t.DB, lastMessage.Message.From.Id)

	_, err = t.TgService.EditMessageReplyMarkup(lastMessage.Message.From.Id, int(lastMsg.TgID), todos)

	return err
}

func (t *TodoController) Delete(lastMessage telegram.MessageWrapper, modelFromCallback telegram.ModelFromCallback) error {
	var todo models.Todo
	todo.Get(t.DB, modelFromCallback.Id)

	if todo.ID == 0 {
		err := errors.New("unknown todo")
		return err
	}

	err := todo.Delete(t.DB, modelFromCallback.Id)

	if err != nil {
		return err
	}

	todos, err := todo.GetAllFromCategoryId(t.DB, todo.CategoryID)

	if err != nil {
		return err
	}

	if todos == nil {
		_, err = t.TgService.EditMessageText(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, "У этой категории больше нет todo")

		if err != nil {
			return err
		}
		return nil
	}

	_, err = t.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, todos)

	return err
}
