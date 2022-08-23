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
	c := make(chan error)
	go func() {
		if user.ID == 0 {
			user.TgID = lastMessage.Message.From.Id
			user.Name = lastMessage.Message.From.FirstName + " " + lastMessage.Message.From.LastName

			err = user.Create(router.DB)
			c <- err
		}
	}()
	go func() {
		err = router.saveIncomingMsg(lastMessage)
		c <- err
	}()
	go func() {
		err = router.saveOutgoingMsg(lastMessage)
		c <- err
	}()

	err = <-c
	if err != nil {
		handleError(w, err)
		return
	}
}
