package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cherryReptile/Todo/internal/models"
	"io"
	"net/http"
	"strconv"
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

// List as inline keyboard button

func (s *Service) SendInlineKeyboard(text string, chatId uint, categories []models.Category) (BotMessage, error) {
	var responseMsg BotMessage

	var inline ToInlineKeyboardBtn
	inline.Text, inline.ChatId = text, chatId
	inline.ReplyMarkup.InlineKeyboard = make([][1]InlineKeyboardBtn, len(categories))

	for i, v := range categories {
		inline.ReplyMarkup.InlineKeyboard[i][0].Text, inline.ReplyMarkup.InlineKeyboard[i][0].CallbackData = v.Name, strconv.Itoa(int(v.ID))
	}

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

// Abstract for doing requests methods

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

//\Here must be reference to struct for response

func (s *Service) AfterRequest(res *http.Response, responseMsg interface{}) error {
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseMsg)

	return err
}
