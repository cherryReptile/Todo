package interfaces

import (
	"github.com/cherryReptile/Todo/internal/database"
	"net/http"
)

type JobInterface interface {
	Init(data interface{})
	Run()
	Close()
	Error(err error)
}

type RequestModelInterface interface {
	CheckBody(r *http.Request) (map[string]string, error)
}

type ModelInterface interface {
	Create(db *database.SqlLite, data map[string]string) error
	Get(db *database.SqlLite, id uint) error
	Update(db *database.SqlLite, id uint) error
	Delete(db *database.SqlLite, id uint) error
}
