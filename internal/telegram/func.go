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
		return updates, errors.New(fmt.Sprintf("telegram not ok"))
	}

	return updates, nil
}

func (s *Service) SendMessage(chatId uint, message string) error {
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
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) HandleMethods(message MessageWrapper) {
	switch message.Message.Text {
	case "/start":
		s.SendHello(message)
		break
	default:
		s.SendDefault(message)
		break
	}
}

func (s *Service) SendHello(message MessageWrapper) {
	msg := fmt.Sprintf("Привет %v! \nДобро пожаловать в наш сервис, наши команды:\n/start - начало\n/list - показать мои todo", message.Message.From.FirstName)
	s.SendMessage(message.Message.From.Id, msg)
}

func (s *Service) SendDefault(message MessageWrapper) {
	msg := fmt.Sprintf("Извините, команда не распознана 😞, наши команды:\n/start - начало\n/list - показать мои todo")
	s.SendMessage(message.Message.From.Id, msg)
}
