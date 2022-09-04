package router

import (
	"github.com/cherryReptile/Todo/internal/controllers"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/telegram"
	"net/http"
)

type Router struct {
	Worker             *queue.JobWorker
	DB                 *database.SqlLite
	TgService          *telegram.Service
	CategoryController *controllers.CategoryController
	TodoController     *controllers.TodoController
}

func NewRouter(Worker *queue.JobWorker, db *database.SqlLite, service *telegram.Service) Router {
	return Router{
		Worker:             Worker,
		DB:                 db,
		TgService:          service,
		CategoryController: controllers.NewCategoryController(db, service),
		TodoController:     controllers.NewTodoController(db, service),
	}
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

	err = router.saveIncomingMsg(lastMessage)

	if err != nil {
		handleError(w, err)
		return
	}

	var lastCommand models.Message
	var callback models.Callback
	var penultimateMsg models.Message

	switch {
	case lastMessage.Message.MessageId != 0:
		err = lastCommand.GetLastCommand(router.DB, lastMessage.Message.From.Id)
		if err != nil {
			break
		}
		err = callback.GetLast(router.DB, lastMessage.Message.From.Id)
		if err != nil {
			break
		}
		err = penultimateMsg.GetLast(router.DB, lastMessage.Message.From.Id)
	case lastMessage.CallbackQuery.Id != "":
		err = lastCommand.GetLastCommand(router.DB, uint(lastMessage.CallbackQuery.Chat.Id))
		if err != nil {
			break
		}
		err = callback.GetLast(router.DB, uint(lastMessage.CallbackQuery.Chat.Id))
		if err != nil {
			break
		}
		err = penultimateMsg.GetLast(router.DB, uint(lastMessage.CallbackQuery.Chat.Id))
	}

	if err != nil {
		handleError(w, err)
		return
	}

	err = router.handleLastCommand(lastCommand, penultimateMsg, callback, lastMessage)

	if err != nil {
		handleError(w, err)
		return
	}
}
