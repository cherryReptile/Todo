package controllers

import (
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

	switch modelFromCallback.Model {
	case "category":
		todos, err := todo.GetAllFromCategoryId(t.DB, category.ID)

		if err != nil {
			break
		}

		if todos == nil {
			_, err = t.TgService.SendMessage(uint(lastMessage.CallbackQuery.Chat.Id), fmt.Sprintf("У %v нет todo, используйте /todoCreate", category.Name))
		} else {
			text := fmt.Sprintf("Todo %v категории\n", category.Name)
			_, err = t.TgService.EditMessageText(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, text)

			if err != nil {
				break
			}

			_, err = t.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, todos)
		}
	default:
		todo.Name = lastMessage.Message.Text
		todo.CategoryID = category.ID
		err := todo.Create(t.DB)

		if err != nil {
			return err
		}

		todos, err := todo.GetAllFromCategoryId(t.DB, category.ID)

		if err != nil {
			return err
		}

		_, err = t.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, todos)
	}

	if err != nil {
		return err
	}

	return nil
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

	_, err = t.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, todos)

	return err
}
