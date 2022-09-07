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

	if category.ID == 0 {
		err := errors.New("unknown category")
		return err
	}

	var todo models.Todo
	todo.Name = lastMessage.Message.Text
	todo.CategoryID = category.ID

	err := todo.Create(t.DB)

	if err != nil {
		return err
	}

	text := fmt.Sprintf("Todo %v создана в категории %v", todo.Name, category.Name)
	botMsg, err := t.TgService.SendMessage(lastMessage.Message.From.Id, text)

	if err != nil {
		return err
	}

	err = SaveBotMsg(botMsg, t.DB)

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
		_, err = t.TgService.EditMessageText(lastMessage.CallbackQuery.From.Id, lastMessage.CallbackQuery.Message.MessageId, "У этой категории больше нет todo")

		if err != nil {
			return err
		}
		return nil
	}

	_, err = t.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.From.Id), lastMessage.CallbackQuery.Message.MessageId, "todoDelete", todos)

	return err
}
