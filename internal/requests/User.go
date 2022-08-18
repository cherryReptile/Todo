package requests

import (
	"encoding/json"
	"github.com/cherryReptile/Todo/internal/validations"
	"net/http"
)

type User struct {
	Request *http.Request
	Name    string `json:"name" validate:"required"`
	TgID    uint   `json:"tg_id" validate:"required"`
}

func NewUser(r *http.Request) User {
	return User{
		Request: r,
	}
}

func (u *User) CheckBody() error {
	err := json.NewDecoder(u.Request.Body).Decode(u)

	if err != nil {
		return err
	}

	err = validations.CreatingOrUpdatingValidate(u)

	return err
}
