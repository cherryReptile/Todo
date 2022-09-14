package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cherryReptile/Todo/internal/models"
	"log"
	"net/http"
	"os"
	"time"
)

// Tg methods

func (s *Service) SetWebhook() error {
	toWh := struct {
		Url string `json:"url"`
	}{}

	toWh.Url = os.Getenv("WEBHOOK_URL")
	res, err := s.DoRequest("setWebhook", "POST", toWh)

	if err != nil {
		return err
	}
	fmt.Println(res)

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("telegram status %v", res.StatusCode))
	}

	return nil
}

func (s *Service) SetCommands() error {
	var commands ToBotCommands
	var command Command

	command.SetFields("/start", "начало")
	commands.Commands[0] = command
	command.SetFields("/category_create", "создать категорию для todo")
	commands.Commands[1] = command
	command.SetFields("/list", "все категории и их todo(которые можно удалить при нажатии)")
	commands.Commands[2] = command
	command.SetFields("/category_delete", "удалить категорию")
	commands.Commands[3] = command
	command.SetFields("/todo_create", "создать todo")
	commands.Commands[4] = command
	command.SetFields("/todo_delete", "удалить todo")
	commands.Commands[5] = command

	res, err := s.DoRequest("setMyCommands", "POST", commands)

	if err != nil {
		return err
	}
	fmt.Println(res)

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("telegram status %v", res.StatusCode))
	}

	if err != nil {
		return err
	}

	return nil
}

//\ for long polling
//
//func (s *Service) GetUpdates() (Updates, error) {
//	var updates Updates
//	var allowedU []string
//
//	toUpdates := ToUpdates{
//		AllowedUpdates: allowedU,
//	}
//
//	res, err := s.DoRequest("getUpdates", "POST", toUpdates)
//
//	if err != nil {
//		return updates, err
//	}
//
//	defer res.Body.Close()
//
//	if res.StatusCode != http.StatusOK {
//		return updates, errors.New(fmt.Sprintf("telegram status %v", res.StatusCode))
//	}
//
//	err = s.AfterRequest(res, &updates)
//
//	if err != nil {
//		return updates, err
//	}
//
//	if !updates.OK {
//		err = errors.New(fmt.Sprintf("telegram not ok"))
//		return updates, err
//	}
//
//	return updates, nil
//}

func (s *Service) SendMessage(chatId uint, message string) (BotMessage, error) {
	var responseMsg BotMessage

	tmsg := ToMessage{
		ChatId: chatId,
		Text:   message,
	}

	res, err := s.DoRequest("sendMessage", "POST", tmsg)

	if err != nil {
		return responseMsg, err
	}

	defer res.Body.Close()

	err = s.AfterRequest(res, &responseMsg)

	if err != nil {
		return responseMsg, err
	}

	if !responseMsg.Ok {
		err = errors.New(fmt.Sprintf("telegram not ok"))
		return responseMsg, err
	}

	return responseMsg, nil
}

//\ this need for not receive overhead callbacks

func (s *Service) AnswerCallbackQuery(callbackId string, text string) error {
	toCallback := ToAnswerCallback{
		CallbackQueryId: callbackId,
		Text:            text,
	}

	res, err := s.DoRequest("answerCallbackQuery", "POST", toCallback)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	var responseMsg BotMessage
	err = s.AfterRequest(res, &responseMsg)

	if err != nil {
		return err
	}

	if !responseMsg.Ok {
		err = errors.New(fmt.Sprintf("telegram not ok"))
		return err
	}

	return nil
}

// List as inline keyboard buttons

func (s *Service) SendInlineKeyboard(text string, chatId uint, btnMethod string, model interface{}) (BotMessage, error) {
	var responseMsg BotMessage

	var inline ToInlineKeyboardBtn
	inline.Text, inline.ChatId = text, chatId

	modelSwitcher(&inline, btnMethod, model)

	res, err := s.DoRequest("sendMessage", "POST", inline)

	if err != nil {
		return responseMsg, err
	}

	defer res.Body.Close()

	err = s.AfterRequest(res, &responseMsg)

	if err != nil {
		return responseMsg, err
	}

	if !responseMsg.Ok {
		err = errors.New(fmt.Sprintf("telegram not ok"))
		return responseMsg, err
	}

	return responseMsg, nil
}

func (s *Service) EditMessageReplyMarkup(chatId uint, messageId int, btnMethod string, model interface{}) (BotMessage, error) {
	var responseMsg BotMessage
	var inline ToInlineKeyboardBtn
	inline.ChatId, inline.MessageId = chatId, messageId

	modelSwitcher(&inline, btnMethod, model)

	res, err := s.DoRequest("editMessageReplyMarkup", "POST", inline)

	if err != nil {
		return responseMsg, err
	}

	defer res.Body.Close()

	err = s.AfterRequest(res, &responseMsg)

	if err != nil {
		return responseMsg, err
	}

	if !responseMsg.Ok {
		err = errors.New(fmt.Sprintf("telegram not ok"))
		return responseMsg, err
	}

	return responseMsg, nil
}

func (s *Service) EditMessageText(chatId uint, messageId int, text string) (BotMessage, error) {
	var responseMsg BotMessage

	tmsg := ToEditMessage{
		chatId, messageId, text,
	}

	res, err := s.DoRequest("editMessageText", "POST", tmsg)

	if err != nil {
		return responseMsg, err
	}

	defer res.Body.Close()

	err = s.AfterRequest(res, &responseMsg)

	if err != nil {
		return responseMsg, err
	}

	if !responseMsg.Ok {
		err = errors.New(fmt.Sprintf("telegram not ok"))
		return responseMsg, err
	}

	return responseMsg, nil
}

func (s *Service) BeforeRequest(tgMethod string, httpMethod string, paramStruct interface{}) (*http.Request, error) {
	url := s.BotUrl + "/" + tgMethod

	req, err := http.NewRequest(httpMethod, url, ToReader(paramStruct))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// Abstract for doing requests methods

func modelSwitcher(inline *ToInlineKeyboardBtn, method string, model interface{}) {
	switch model.(type) {
	case []models.Category:
		categories := model.([]models.Category)
		inline.ReplyMarkup.InlineKeyboard = make([][1]InlineKeyboardBtn, len(categories))

		for i, v := range categories {
			var category ModelFromCallback
			category.Id = v.ID
			category.Model = "category"
			category.Method = method

			jsonBytes, err := json.Marshal(category)

			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Println(string(jsonBytes))
			inline.ReplyMarkup.InlineKeyboard[i][0].Text, inline.ReplyMarkup.InlineKeyboard[i][0].CallbackData = v.Name, string(jsonBytes)
		}
	case []models.Todo:
		todos := model.([]models.Todo)
		inline.ReplyMarkup.InlineKeyboard = make([][1]InlineKeyboardBtn, len(todos))

		for i, v := range todos {
			var todo ModelFromCallback
			todo.Id = v.ID
			todo.Model = "todo"
			todo.Method = method

			jsonBytes, err := json.Marshal(todo)

			if err != nil {
				log.Fatal(err)
				return
			}
			inline.ReplyMarkup.InlineKeyboard[i][0].Text, inline.ReplyMarkup.InlineKeyboard[i][0].CallbackData = v.Name, string(jsonBytes)
		}
	}
}

func (s *Service) DoRequest(tgMethod string, httpMethod string, paramStruct interface{}) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 20,
	}

	req, err := s.BeforeRequest(tgMethod, httpMethod, paramStruct)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

//\Here must be a reference to struct for response

func (s *Service) AfterRequest(res *http.Response, responseMsg interface{}) error {
	err := json.NewDecoder(res.Body).Decode(responseMsg)

	return err
}
