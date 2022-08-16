package router

import (
	"encoding/json"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/jobs"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/queue"
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
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		handleError(w, err)
		return
	}

	db := router.DB.DB
	result := db.Select("Name", "TgID").Create(&u)

	if result.Error != nil {
		handleError(w, result.Error)
		return
	}

	db.First(&u, u.ID)

	responseJson(w, u)
}

func (router *Router) UserGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		handleError(w, err)
		return
	}

	u := new(models.User)
	result := router.DB.DB.First(&u, id)

	if result.Error != nil {
		handleError(w, result.Error)
		return
	}

	responseJson(w, u)
}

func (router *Router) UserUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		handleError(w, err)
		return
	}

	u := new(models.User)
	err = json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		handleError(w, err)
		return
	}

	name := u.Name
	result := router.DB.DB.First(&u, id)

	if result.Error != nil {
		handleError(w, result.Error)
		return
	}

	result = router.DB.DB.Model(&u).Update("name", name)

	if result.Error != nil {
		handleError(w, result.Error)
		return
	}

	responseJson(w, u)
}

func (router *Router) UserDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		handleError(w, err)
		return
	}

	u := new(models.User)
	result := router.DB.DB.First(&u, id)

	if result.Error != nil {
		handleError(w, result.Error)
		return
	}

	router.DB.DB.Delete(&u)

	w.WriteHeader(204)
}
