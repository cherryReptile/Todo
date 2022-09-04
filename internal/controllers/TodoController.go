package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/telegram"
	"strconv"
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

func (t *TodoController) Create(lastMessage telegram.MessageWrapper, callbackMsg models.Callback) error {
	go AnswerToCallback(lastMessage, t.TgService)
	var user models.User
	user.GetFromTg(t.DB, uint(lastMessage.CallbackQuery.Message.Chat.Id))

	var callback telegram.CallbackQuery
	json.Unmarshal([]byte(callbackMsg.Json), &callback)

	var category models.Category
	id, err := strconv.Atoi(callback.Data)

	if err != nil {
		return err
	}

	category.Get(t.DB, uint(id))

	var todo models.Todo
	todos, err := todo.GetAllFromCategoryId(t.DB, category.ID)

	if err != nil {
		return err
	}

	if todos == nil {
		_, err = t.TgService.SendMessage(uint(lastMessage.CallbackQuery.Chat.Id), fmt.Sprintf("У %v нет todo", category.Name))
	} else if lastMessage.Message.MessageId != 0 {
		todo.Name = lastMessage.Message.Text
		todo.CategoryID = category.ID
		err = todo.Create(t.DB)

		if err != nil {
			return err
		}

		todos, err = todo.GetAllFromCategoryId(t.DB, category.ID)

		if err != nil {
			return err
		}

		text := fmt.Sprintf("Todo %v категории\n", category.Name)

		for i, v := range todos {
			text += fmt.Sprintf("%v.%v\n", i+1, v.Name)
		}

		_, err = t.TgService.EditMessageText(lastMessage.Message.From.Id, callback.Message.MessageId, text)
	} else if lastMessage.CallbackQuery.Id != "" {
		text := fmt.Sprintf("Todo %v категории\n", category.Name)

		for i, v := range todos {
			text += fmt.Sprintf("%v.%v\n", i+1, v.Name)
		}

		if lastMessage.CallbackQuery.Message.Text != text {
			_, err = t.TgService.EditMessageText(uint(lastMessage.CallbackQuery.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, text)
		}

		if err != nil {
			return err
		}

		_, err = t.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, []models.Todo{})

	}

	if err != nil {
		return err
	}

	return nil
}
