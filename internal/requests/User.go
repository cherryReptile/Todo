package requests

import (
	"net/http"
)

type User struct {
	Name string `json:"name" validate:"required"`
	TgID uint   `json:"tg_id" validate:"required"`
}

func (u *User) CheckBody(r *http.Request) (map[string]string, error) {
	err := checkBody(r, u)

	if err == nil {
		data := make(map[string]string)
		data["name"] = u.Name
		data["tg_id"] = string(u.TgID)
		return data, nil
	}

	return nil, err
}
