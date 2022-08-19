package requests

import (
	"encoding/json"
	"github.com/cherryReptile/Todo/internal/validations"
	"net/http"
)

func checkBody(r *http.Request, i interface{}) error {
	err := json.NewDecoder(r.Body).Decode(i)

	if err != nil {
		return err
	}

	err = validations.CreatingOrUpdatingValidate(i)

	return err
}
