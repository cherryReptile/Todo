package router

import (
	"encoding/json"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/jobs"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/responses"
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

func (router *Router) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	ResponseError(w, err)

	err = u.Create(router.DB)
	ResponseError(w, err)

	responseJson(w, u)
}

func (router *Router) GetUser(w http.ResponseWriter, r *http.Request) {
	u := new(models.User)
	err := u.Get(router.DB, 1)
	ResponseError(w, err)

	responseJson(w, u)
}

//func (router *Router) CreateCategory(w http.ResponseWriter, r *http.Request) {
//	var c models.Category
//	err := json.NewDecoder(r.Body).Decode(&c)
//
//	if err != nil {
//		handleError(w, err)
//		return
//	}
//	err = c.Create(router.DB)
//
//	if err != nil {
//		handleError(w, err)
//		return
//	}
//	responseJson(w, c)
//}
