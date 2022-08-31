package telegram

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/models"
)

//Messages layouts

func (s *Service) HandleMethods(message MessageWrapper) (BotMessage, error) {
	var botMsg BotMessage
	var err error

	switch message.Message.Text {
	case "/start":
		botMsg, err = s.SendHello(message)
		break
	case "/categoryCreate":
		botMsg, err = s.SendCreate(message)
		break
	//case "/categoryDelete":
	//	botMsg, err = s.SendInlineKeyboard("Выберите какую категорию с её todo удалить:", message.Message.From.Id)
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

func (s *Service) SendList(message MessageWrapper, categories []models.Category) (BotMessage, error) {
	var msg string
	msg = "Твои категории: 👇\n"
	for _, v := range categories {
		msg += fmt.Sprintf("%v\n", v.Name)
	}
	return s.SendMessage(message.Message.From.Id, msg)
}
