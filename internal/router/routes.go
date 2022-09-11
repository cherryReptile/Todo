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

func NewRouter(Worker *queue.JobWorker, db *database.SqlLite, service *telegram.Service) *Router {
	return &Router{
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

	var modelFromCallback telegram.ModelFromCallback
	lastCommand, err := router.toCallback(lastMessage, &modelFromCallback)

	if err != nil {
		handleError(w, err)
		return
	}

	err = router.handleLastCommand(lastCommand, modelFromCallback, lastMessage)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, map[string]string{
		"status": "ok",
	})
}
