package router

import (
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/jobs"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/requests"
	"github.com/cherryReptile/Todo/internal/responses"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func (router *Router) UserCreate(w http.ResponseWriter, r *http.Request) {
	reqU := requests.NewUser(r)
	err := reqU.CheckBody()

	if err != nil {
		handleError(w, err)
		return
	}

	u := new(models.User)
	err = u.Create(router.DB, reqU)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, u)
}

func (router *Router) UserGet(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	u := new(models.User)
	err = u.Get(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, u)
}

func (router *Router) UserUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	reqU := requests.NewUser(r)
	err = reqU.CheckBody()

	if err != nil {
		handleError(w, err)
		return
	}

	u := new(models.User)
	u.Name = reqU.Name
	err = u.Update(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, u)
}

func (router *Router) UserDelete(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	err = new(models.User).Delete(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(204)
}

func convertId(key string, r *http.Request) (uint, error) {
	id, err := strconv.Atoi(mux.Vars(r)[key])
	return uint(id), err
}
