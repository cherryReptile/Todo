package router

import (
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/responses"
	"github.com/cherryReptile/Todo/internal/telegram"
	"net/http"
	"strconv"
)

type Router struct {
	Worker    *queue.JobWorker
	DB        *database.SqlLite
	TgService *telegram.Service
}

func NewRouter(Worker *queue.JobWorker, db *database.SqlLite, service *telegram.Service) Router {
	return Router{
		Worker:    Worker,
		DB:        db,
		TgService: service,
	}
}

func (router *Router) Index(w http.ResponseWriter, r *http.Request) {
	responseJson(w, responses.VersionResponse{Version: "1.0", Name: "Starter Kit v1.0"})
}

func (router *Router) Start(w http.ResponseWriter, r *http.Request) {
	lastMessage, err := router.getLastMsg()

	if err != nil {
		handleError(w, err)
		return
	}

	var user models.User

	user.GetFromTg(router.DB, lastMessage.Message.From.Id)
	go func() {
		if user.ID == 0 {
			user.TgID = lastMessage.Message.From.Id
			user.Name = lastMessage.Message.From.FirstName + " " + lastMessage.Message.From.LastName

			err = user.Create(router.DB)
			if err != nil {
				handleError(w, err)
				return
			}
		}
	}()
	go func() {
		err = router.saveIncomingMsg(lastMessage)
		if err != nil {
			handleError(w, err)
			return
		}
	}()

	err = router.saveHandledMsg(lastMessage)

	if err != nil {
		handleError(w, err)
		return
	}
}

func (router Router) CategoryCreate(w http.ResponseWriter, r *http.Request) {
	lastMessage, err := router.getLastMsg()

	if err != nil {
		handleError(w, err)
		return
	}

	err = router.saveIncomingMsg(lastMessage)
	if err != nil {
		handleError(w, err)
	}

	var msg models.Message
	msgs, err := msg.GetLastTwo(router.DB, lastMessage.Message.From.Id)

	if err != nil {
		handleError(w, err)
		return
	}

	if msgs[1].Command.Valid && !msgs[0].Command.Valid {
		router.handleLastCommand(msgs[1], lastMessage)
		return
	}

	err = router.saveHandledMsg(lastMessage)

	if err != nil {
		handleError(w, err)
		return
	}
}

func (router *Router) CategoryList(w http.ResponseWriter, r *http.Request) {
	lastMessage, err := router.getLastMsg()

	if err != nil {
		handleError(w, err)
		return
	}

	go func() {
		err = router.saveIncomingMsg(lastMessage)
		if err != nil {
			handleError(w, err)
		}
	}()

	var user models.User
	user.GetFromTg(router.DB, lastMessage.Message.From.Id)

	var category models.Category
	categories, err := category.GetAllCategories(router.DB, user.ID)

	if err != nil {
		handleError(w, err)
		return
	}

	botMsg, err := router.TgService.SendInlineKeyboard("Твои категории: 👇\n", lastMessage.Message.From.Id, categories)

	if err != nil {
		handleError(w, err)
		return
	}

	err = router.saveBotMsg(botMsg)

	if err != nil {
		handleError(w, err)
		return
	}
}

func (router *Router) CategoryGet(w http.ResponseWriter, r *http.Request) {
	lastUpdate, err := router.getLastMsg()

	if err != nil {
		handleError(w, err)
		return
	}

	go router.AnswerToCallback(lastUpdate)

	var category models.Category
	id, err := strconv.Atoi(lastUpdate.CallbackQuery.Data)

	if err != nil {
		handleError(w, err)
		return
	}

	category.Get(router.DB, uint(id))

	botMsg, err := router.TgService.SendMessage(uint(lastUpdate.CallbackQuery.Chat.Id), category.Name)

	if err != nil {
		handleError(w, err)
		return
	}

	err = router.saveBotMsg(botMsg)

	if err != nil {
		handleError(w, err)
		return
	}
}

func (router Router) CategoryDelete(w http.ResponseWriter, r *http.Request) {
	lastUpdate, err := router.getLastMsg()

	if err != nil {
		handleError(w, err)
		return
	}

	if lastUpdate.Message.Text == "/categoryDelete" {
		go func() {
			if lastUpdate.CallbackQuery.Id == "" {
				err = router.saveIncomingMsg(lastUpdate)
				if err != nil {
					handleError(w, err)
				}
			}
			return
		}()

		var user models.User
		user.GetFromTg(router.DB, lastUpdate.Message.From.Id)
		var category models.Category
		categories, err := category.GetAllCategories(router.DB, user.ID)

		if err != nil {
			handleError(w, err)
			return
		}

		botMsg, err := router.TgService.SendInlineKeyboard("Выберите какую категорию удалить\n", lastUpdate.Message.From.Id, categories)

		if err != nil {
			handleError(w, err)
			return
		}

		err = router.saveBotMsg(botMsg)
	} else if lastUpdate.CallbackQuery.Id != "" {
		go router.AnswerToCallback(lastUpdate)
		var user models.User
		user.GetFromTg(router.DB, uint(lastUpdate.CallbackQuery.Message.Chat.Id))

		var category models.Category
		id, err := strconv.Atoi(lastUpdate.CallbackQuery.Data)

		if err != nil {
			handleError(w, err)
			return
		}

		category.Get(router.DB, uint(id))

		category.Delete(router.DB, uint(id))

		categories, err := category.GetAllCategories(router.DB, user.ID)

		if err != nil {
			handleError(w, err)
			return
		}

		_, err = router.TgService.EditMessageReplyMarkup(uint(lastUpdate.CallbackQuery.Message.Chat.Id), lastUpdate.CallbackQuery.Message.MessageId, categories)

		if err != nil {
			handleError(w, err)
			return
		}
	}
}
