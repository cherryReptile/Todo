package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cherryReptile/Todo/internal/models"
	"io"
	"log"
	"net/http"
	"time"
)

// Tg metgods

func (s *Service) GetUpdates() (Updates, error) {
	var updates Updates
	var allowedU []string

	toUpdates := ToUpdates{
		AllowedUpdates: allowedU,
	}

	res, err := s.DoRequest("getUpdates", "POST", toUpdates)

	if err != nil {
		return updates, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return updates, errors.New(fmt.Sprintf("telegram status %v", res.StatusCode))
	}

	err = s.AfterRequest(res, &updates)

	if err != nil {
		return updates, err
	}

	if !updates.OK {
		err = errors.New(fmt.Sprintf("telegram not ok"))
		return updates, err
	}

	return updates, nil
}

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

func (s *Service) SendInlineKeyboard(text string, chatId uint, model interface{}) (BotMessage, error) {
	var responseMsg BotMessage

	var inline ToInlineKeyboardBtn
	inline.Text, inline.ChatId = text, chatId

	modelSwitcher(&inline, model)

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

func (s *Service) EditMessageReplyMarkup(chatId uint, messageId int, model interface{}) (BotMessage, error) {
	var responseMsg BotMessage
	var inline ToInlineKeyboardBtn
	inline.ChatId, inline.MessageId = chatId, messageId

	modelSwitcher(&inline, model)

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

func modelSwitcher(inline *ToInlineKeyboardBtn, model interface{}) {
	switch model.(type) {
	case []models.Category:
		categories := model.([]models.Category)
		inline.ReplyMarkup.InlineKeyboard = make([][1]InlineKeyboardBtn, len(categories))

		for i, v := range categories {
			var category ModelFromCallback
			category.Id = v.ID
			category.Model = "category"

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
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseMsg)

	return err
}
