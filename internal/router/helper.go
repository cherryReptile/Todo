package router

import (
	"encoding/json"
	"github.com/cherryReptile/Todo/internal/responses"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func responseJson(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(r)

	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadGateway)

	response := responses.ErrorResponse{
		Errors: []string{err.Error()},
	}

	r, _ := json.Marshal(response)

	w.Write(r)
}

func convertId(key string, r *http.Request) (uint, error) {
	id, err := strconv.Atoi(mux.Vars(r)[key])
	return uint(id), err
}
