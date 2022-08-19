package requests

import (
	"net/http"
)

type User struct {
	Name string `json:"name" validate:"required"`
	TgID uint   `json:"tg_id" validate:"required"`
}

func (u *User) CheckBody(r *http.Request) error {
	err := checkBody(r, u)

	return err
}
