package router

import (
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/telegram"
	"strconv"
)

func (router Router) Create(lastMessage telegram.MessageWrapper) error {
	var user models.User
	user.GetFromTg(router.DB, lastMessage.Message.From.Id)

	var category models.Category
	category.Name = lastMessage.Message.Text
	category.UserID = user.ID

	err := category.Create(router.DB)

	if err != nil {
		return err
	}

	botMsg, err := router.TgService.SendCreated(lastMessage)

	if err != nil {
		return err
	}

	err = router.saveBotMsg(botMsg)

	return err
}

func (router *Router) List(lastMessage telegram.MessageWrapper, text string) error {
	var user models.User
	user.GetFromTg(router.DB, lastMessage.Message.From.Id)

	var category models.Category
	categories, err := category.GetAllCategories(router.DB, user.ID)

	if err != nil {
		return err
	}

	botMsg, err := router.TgService.SendInlineKeyboard(text, lastMessage.Message.From.Id, categories)

	if err != nil {
		return err
	}

	err = router.saveBotMsg(botMsg)

	if err != nil {
		return err
	}

	return nil
}

func (router *Router) Get(lastMessage telegram.MessageWrapper) error {
	go router.AnswerToCallback(lastMessage)

	var category models.Category
	id, err := strconv.Atoi(lastMessage.CallbackQuery.Data)

	if err != nil {
		return err
	}

	category.Get(router.DB, uint(id))

	botMsg, err := router.TgService.SendMessage(uint(lastMessage.CallbackQuery.Chat.Id), category.Name)

	if err != nil {
		return err
	}

	err = router.saveBotMsg(botMsg)

	if err != nil {
		return err
	}

	return nil
}

func (router *Router) Delete(lastMessage telegram.MessageWrapper) error {
	go router.AnswerToCallback(lastMessage)
	var user models.User
	user.GetFromTg(router.DB, uint(lastMessage.CallbackQuery.Message.Chat.Id))

	var category models.Category
	id, err := strconv.Atoi(lastMessage.CallbackQuery.Data)

	if err != nil {
		return err
	}

	category.Delete(router.DB, uint(id))

	categories, err := category.GetAllCategories(router.DB, user.ID)

	if err != nil {
		return err
	}

	_, err = router.TgService.EditMessageReplyMarkup(uint(lastMessage.CallbackQuery.Message.Chat.Id), lastMessage.CallbackQuery.Message.MessageId, categories)

	if err != nil {
		return err
	}

	return nil
}
