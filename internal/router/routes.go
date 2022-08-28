package router

import (
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/responses"
	"github.com/cherryReptile/Todo/internal/telegram"
	"net/http"
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

	botMsg, err := router.TgService.SendList(lastMessage, categories)

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
