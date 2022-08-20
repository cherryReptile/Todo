package requests

import "net/http"

type Todo struct {
	Name string `json:"name" validate:"max=255,min=1"`
}

func (t *Todo) CheckBody(r *http.Request) error {
	err := checkBody(r, t)

	return err
}
