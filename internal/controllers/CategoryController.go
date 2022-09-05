package controllers

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/telegram"
	"strconv"
)

type CategoryController struct {
	DB        *database.SqlLite
	TgService *telegram.Service
}

func NewCategoryController(db *database.SqlLite, service *telegram.Service) *CategoryController {
	return &CategoryController{
		DB:        db,
		TgService: service,
	}
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

func (c *CategoryController) List(lastMessage telegram.MessageWrapper, text string) error {
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
		botMsg, err = c.TgService.SendInlineKeyboard(text, lastMessage.Message.From.Id, categories)
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
		botMsg, err = c.TgService.SendMessage(uint(lastMessage.CallbackQuery.Chat.Id), fmt.Sprintf("У %v нет todo, используйте /todo, чтобы создать", category.Name))
	} else {
		botMsg, err = c.TgService.SendInlineKeyboard(fmt.Sprintf("Todo %v категории(нажми, чтобы удалить):", category.Name), uint(lastMessage.CallbackQuery.Chat.Id), todos)
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

func (c *CategoryController) Delete(lastMessage telegram.MessageWrapper) error {
	go AnswerToCallback(lastMessage, c.TgService)
	var user models.User
	user.GetFromTg(c.DB, uint(lastMessage.CallbackQuery.Message.Chat.Id))

	var category models.Category
	id, err := strconv.Atoi(lastMessage.CallbackQuery.Data)

	if err != nil {
		return err
	}

	category.Delete(c.DB, uint(id))

	categories, err := category.GetAllCategories(c.DB, user.ID)

	if err != nil {
		return err
	}

	_, err = c.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, categories)

	if err != nil {
		return err
	}

	return nil
}
