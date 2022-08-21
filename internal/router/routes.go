package router

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/controllers"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/jobs"
	"github.com/cherryReptile/Todo/internal/models"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/requests"
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

func (router *Router) UserCreate(w http.ResponseWriter, r *http.Request) {
	controller := controllers.NewController(router.DB, r)
	reqU := new(requests.User)
	u := new(models.User)
	err := controller.AbstractCreate(reqU, u)
	fmt.Println(err)
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

//func (router *Router) UserUpdate(w http.ResponseWriter, r *http.Request) {
//	id, err := convertId("id", r)
//
//	if err != nil {
//		handleError(w, err)
//		return
//	}
//
//	reqU := new(requests.User)
//	err = reqU.CheckBody(r)
//
//	if err != nil {
//		handleError(w, err)
//		return
//	}
//
//	u := new(models.User)
//	u.Name = reqU.Name
//	err = u.Update(router.DB, id)
//
//	if err != nil {
//		handleError(w, err)
//		return
//	}
//
//	responseJson(w, u)
//}

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

func (router *Router) CategoryCreate(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("user_id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	reqC := new(requests.Category)
	err = reqC.CheckBody(r)

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

	c := new(models.Category)
	c.UserID = u.ID
	err = c.Create(router.DB, reqC)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, c)
}

func (router *Router) CategoryGet(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	c := new(models.Category)
	err = c.Get(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, c)
}

func (router *Router) CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	reqC := new(requests.Category)
	err = reqC.CheckBody(r)

	if err != nil {
		handleError(w, err)
		return
	}

	c := new(models.Category)
	c.Name = reqC.Name
	err = c.Update(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, c)
}

func (router *Router) CategoryDelete(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	err = new(models.Category).Delete(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(204)
}

func (router *Router) TodoCreate(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("category_id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	reqT := new(requests.Todo)
	err = reqT.CheckBody(r)

	if err != nil {
		handleError(w, err)
		return
	}

	c := new(models.Category)
	err = c.Get(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	t := new(models.Todo)
	t.CategoryID = c.ID
	err = t.Create(router.DB, reqT)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, t)
}

func (router *Router) TodoGet(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	t := new(models.Todo)
	err = t.Get(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, t)
}

func (router *Router) TodoUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	reqT := new(requests.Todo)
	err = reqT.CheckBody(r)

	if err != nil {
		handleError(w, err)
		return
	}

	t := new(models.Todo)
	t.Name = reqT.Name
	err = t.Update(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	responseJson(w, t)
}

func (router *Router) TodoDelete(w http.ResponseWriter, r *http.Request) {
	id, err := convertId("id", r)

	if err != nil {
		handleError(w, err)
		return
	}

	err = new(models.Todo).Delete(router.DB, id)

	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(204)
}

//func AbstractCreate(req interfaces.RequestModelInterface, r *http.Request, db *database.SqlLite) (interfaces.ModelInterface, error) {
//	u, err := AbstractCheck(req, r)
//
//	if err != nil {
//		return nil, err
//	}
//
//	err = u.Create(db)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return u, err
//}
//
//func AbstractCheck(req interfaces.RequestModelInterface, r *http.Request) (interfaces.ModelInterface, error) {
//	u, err := req.CheckBody(r)
//	return u, err
//}

//func AbstractGet(model interfaces.ModelInterface, r *http.Request, db *database.SqlLite) error {
//	id, err := convertId("id", r)
//
//	if err != nil {
//		return err
//	}
//
//	err = model.Get(db, id)
//
//	return err
//}

//func (router Router) AbstractUpdate(model interfaces.ModelInterface, r *http.Request) error {
//	id, err := convertId("id", r)
//
//	if err != nil {
//		return err
//	}
//
//	req := model.ItsRequest()
//	err = AbstractCheck(req, r)
//
//	if err != nil {
//		return err
//	}
//
//	u := new(models.User)
//	u.Name = reqU.Name
//	err = u.Update(router.DB, id)
//
//	if err != nil {
//		handleError(w, err)
//		return
//	}
//}
