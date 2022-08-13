package router

import (
	"encoding/json"
	"github.com/pavel-one/GoStarter/internal/database"
	"github.com/pavel-one/GoStarter/internal/interfaces"
	"github.com/pavel-one/GoStarter/internal/jobs"
	"github.com/pavel-one/GoStarter/internal/models"
	"github.com/pavel-one/GoStarter/internal/queue"
	"github.com/pavel-one/GoStarter/internal/responses"
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

func (router *Router) Test(w http.ResponseWriter, r *http.Request) {
	var j jobs.TestJob
	var i interface{}

	j.Init(i)

	router.Worker.Add(&j)
}

func (router *Router) Create(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		handleError(w, err)
		return
	}
	err = u.Create(router.DB)

	if err != nil {
		handleError(w, err)
		return
	}
	responseJson(w, u.Name)
}

func checkInterface(model interfaces.ModelInterface) {

}
