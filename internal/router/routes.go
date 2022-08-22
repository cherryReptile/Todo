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
	Worker *queue.JobWorker
	DB     *database.SqlLite
}

func NewRouter(Worker *queue.JobWorker, db *database.SqlLite) Router {
	return Router{
		Worker: Worker,
		DB:     db,
	}
}

func (router *Router) Index(w http.ResponseWriter, r *http.Request) {
	responseJson(w, responses.VersionResponse{Version: "1.0", Name: "Starter Kit v1.0"})
}

func (router *Router) Start(w http.ResponseWriter, r *http.Request) {
	tgs := new(telegram.Service)
	tgs.Init(router.DB)

	updates, err := tgs.GetUpdates()

	if err != nil {
		handleError(w, err)
		return
	}

	lastMessage := updates.Result[len(updates.Result)-1]

	go func() {
		var user models.User

		user.GetFromTg(router.DB, lastMessage.Message.From.Id)

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

	go tgs.HandleMethods(lastMessage)
}

func (router *Router) CategoryCreate(w http.ResponseWriter, r *http.Request) {
	tgs := new(telegram.Service)
	tgs.Init(router.DB)

	updates, err := tgs.GetUpdates()

	if err != nil {
		handleError(w, err)
		return
	}

	lastMessage := updates.Result[len(updates.Result)-1]
	tgs.HandleMethods(lastMessage)

	for lastMessage.Message.Text == "/createCategory" {
		updates, err = tgs.GetUpdates()
		lastMessage = updates.Result[len(updates.Result)-1]

		if err != nil {
			break
		}
	}

	if err != nil {
		handleError(w, err)
		return
	}

	var user models.User
	user.GetFromTg(router.DB, lastMessage.Message.From.Id)

	var category models.Category
	category.Name = lastMessage.Message.Text
	category.UserID = user.ID

	err = category.Create(router.DB)

	if err != nil {
		responseJson(w, r)
		return
	}

	go tgs.SendCreated(lastMessage)
}
