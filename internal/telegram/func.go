package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (s *Service) GetUpdates() (Updates, error) {
	var updates Updates

	url := s.BotUrl + "/getUpdates"
	method := "GET"

	client := &http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return updates, err
	}

	res, err := client.Do(req)
	if err != nil {
		return updates, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return updates, errors.New(fmt.Sprintf("telegram status %v", res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return updates, err
	}

	err = json.Unmarshal(body, &updates)

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

	url := s.BotUrl + "/sendMessage"
	method := "POST"

	tmsg := ToMessage{
		ChatId: chatId,
		Text:   message,
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest(method, url, tmsg.ToReader())

	if err != nil {
		return responseMsg, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return responseMsg, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return responseMsg, err
	}

	err = json.Unmarshal(body, &responseMsg)

	if err != nil {
		return responseMsg, err
	}

	if !responseMsg.Ok {
		err = errors.New(fmt.Sprintf("telegram not ok"))
		return responseMsg, err
	}

	return responseMsg, nil
}

func (s *Service) HandleMethods(message MessageWrapper) (BotMessage, error) {
	var botMsg BotMessage
	var err error

	switch message.Message.Text {
	case "/start":
		botMsg, err = s.SendHello(message)
		break
	case "/createCategory":
		botMsg, err = s.SendCreate(message)
		break
	default:
		botMsg, err = s.SendDefault(message)
		break
	}
	return botMsg, err
}

func (s *Service) SendHello(message MessageWrapper) (BotMessage, error) {
	msg := fmt.Sprintf("Привет %v! \nДобро пожаловать в наш сервис, наши команды:\n/start - начало\n/list - показать мои todo\n/categoryCreate - создать категорию", message.Message.From.FirstName)
	return s.SendMessage(message.Message.From.Id, msg)
}

func (s *Service) SendDefault(message MessageWrapper) (BotMessage, error) {
	msg := fmt.Sprintf("Извините, команда не распознана 😞, наши команды:\n/start - начало\n/list - показать мои todo")
	return s.SendMessage(message.Message.From.Id, msg)
}

func (s *Service) SendCreate(message MessageWrapper) (BotMessage, error) {
	msg := fmt.Sprintf("Введите название категории")
	return s.SendMessage(message.Message.From.Id, msg)
}

func (s *Service) SendCreated(message MessageWrapper) (BotMessage, error) {
	msg := fmt.Sprintf("Категория %v создана", message.Message.Text)
	return s.SendMessage(message.Message.From.Id, msg)
}
