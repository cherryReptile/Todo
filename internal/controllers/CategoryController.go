package controllers

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/telegram"
)

type CategoryController struct {
	DbController
	TgController
}

func NewCategoryController(db *database.SqlLite, service *telegram.Service) *CategoryController {
	c := new(CategoryController)
	c.DB = db
	c.TgService = service
	return c
}

func (c *CategoryController) Create(lastMessage telegram.MessageWrapper) error {
	var user models.User
	user.GetFromTg(c.DB, lastMessage.Message.From.Id)

	var category models.Category
	category.Name = lastMessage.Message.Text
	category.UserID = user.ID

	err := category.Create(c.DB)

	if err != nil {
		return err
	}

	botMsg, err := c.TgService.SendCreated(lastMessage)

	if err != nil {
		return err
	}

	err = SaveBotMsg(botMsg, c.DB)

	return err
}

func (c *CategoryController) List(lastMessage telegram.MessageWrapper, text string, btnMethod string) error {
	var user models.User
	user.GetFromTg(c.DB, lastMessage.Message.From.Id)

	var category models.Category
	categories, err := category.GetAllCategories(c.DB, user.ID)

	if err != nil {
		return err
	}

	var botMsg telegram.BotMessage

	if categories == nil {
		botMsg, err = c.TgService.SendMessage(lastMessage.Message.From.Id, "У вас ещё нет категорий, используйте 👉 /categoryCreate")
	} else {
		botMsg, err = c.TgService.SendInlineKeyboard(text, lastMessage.Message.From.Id, btnMethod, categories)
	}

	if err != nil {
		return err
	}

	err = SaveBotMsg(botMsg, c.DB)

	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryController) Get(lastMessage telegram.MessageWrapper, modelFromCallback telegram.ModelFromCallback) error {
	go AnswerToCallback(lastMessage, c.TgService)

	var category models.Category

	category.Get(c.DB, modelFromCallback.Id)

	var todo models.Todo
	todos, err := todo.GetAllFromCategoryId(c.DB, category.ID)

	var botMsg telegram.BotMessage

	if todos == nil {
		botMsg, err = c.TgService.EditMessageText(uint(lastMessage.CallbackQuery.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, fmt.Sprintf("У %v нет todo, используйте /todoCreate, чтобы создать", category.Name))
	} else {
		botMsg, err = c.TgService.EditMessageText(uint(lastMessage.CallbackQuery.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, fmt.Sprintf("Todo %v категории(нажми, чтобы удалить):", category.Name))
		botMsg, err = c.TgService.EditMessageReplyMarkup(lastMessage.CallbackQuery.From.Id, lastMessage.CallbackQuery.Message.MessageId, "todoDelete", todos)
	}

	if err != nil {
		return err
	}

	err = SaveBotMsg(botMsg, c.DB)

	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryController) Delete(lastMessage telegram.MessageWrapper, modelFromCallback telegram.ModelFromCallback) error {
	go AnswerToCallback(lastMessage, c.TgService)
	var user models.User
	user.GetFromTg(c.DB, uint(lastMessage.CallbackQuery.Message.Chat.Id))

	var category models.Category
	category.Delete(c.DB, modelFromCallback.Id)

	categories, err := category.GetAllCategories(c.DB, user.ID)

	if err != nil {
		return err
	}

	if categories == nil {
		_, err = c.TgService.EditMessageText(lastMessage.CallbackQuery.From.Id, lastMessage.CallbackQuery.Message.MessageId, "У вас больше нет категорий")
	} else {
		_, err = c.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, "categoryDelete", categories)
	}

	if err != nil {
		return err
	}

	return nil
}
