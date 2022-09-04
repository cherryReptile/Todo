package telegram

import (
	"fmt"
)

//Messages layouts

func (s *Service) SendHello(message MessageWrapper) (BotMessage, error) {
	msg := fmt.Sprintf("Привет %v! \nДобро пожаловать в наш сервис, наши команды:\n/start - начало\n/list - показать мои todo\n/categoryCreate - создать категорию", message.Message.From.FirstName)
	return s.SendMessage(message.Message.From.Id, msg)
}

func (s *Service) SendDefault(message MessageWrapper) (BotMessage, error) {
	msg := "Извините, команда не распознана 😞, наши команды:\n/start - начало\n/categoryCreate - создать категорию для todo\n/list - все категории"
	return s.SendMessage(message.Message.From.Id, msg)
}

func (s *Service) SendCreate(message MessageWrapper) (BotMessage, error) {
	msg := "Введите название категории"
	return s.SendMessage(message.Message.From.Id, msg)
}

func (s *Service) SendCreated(message MessageWrapper) (BotMessage, error) {
	msg := fmt.Sprintf("Категория %v создана", message.Message.Text)
	return s.SendMessage(message.Message.From.Id, msg)
}
