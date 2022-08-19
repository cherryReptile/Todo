package requests

import (
	"net/http"
)

type Category struct {
	Name string `json:"name" validate:"max:255,min:1"`
}

func (c *Category) CheckBody(r *http.Request) error {
	err := checkBody(r, c)

	return err
}
