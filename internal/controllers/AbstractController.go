package controllers

import (
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/interfaces"
	"net/http"
)

type Controller struct {
	DB      *database.SqlLite
	Request *http.Request
}

func NewController(db *database.SqlLite, r *http.Request) Controller {
	return Controller{
		DB:      db,
		Request: r,
	}
}
func (c *Controller) AbstractCreate(req interfaces.RequestModelInterface, model interfaces.ModelInterface) error {

	data, err := AbstractCheck(req, c.Request)

	if err != nil {
		return err
	}

	err = model.Create(c.DB, data)

	return err
}

func AbstractCheck(req interfaces.RequestModelInterface, r *http.Request) (map[string]string, error) {
	data, err := req.CheckBody(r)
	return data, err
}

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
