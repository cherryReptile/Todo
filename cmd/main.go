package main

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/router"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	q := queue.Run()
	sql := database.Connect()
	defer sql.DB.Close()
	route := router.NewRouter(&q, &sql)
	r := mux.NewRouter()
	s := r.Host("127.0.0.1:3000").Subrouter()

	s.HandleFunc("/", route.Index).Methods("GET")
	s.HandleFunc("/job", route.Test)
	s.HandleFunc("/user", route.UserCreate).Methods("POST")
	s.HandleFunc("/user/{id}", route.UserGet).Methods("GET")
	s.HandleFunc("/user/{id}", route.UserUpdate).Methods("PATCH")
	s.HandleFunc("/user/{id}", route.UserDelete).Methods("DELETE")
	s.HandleFunc("/{user_id}/category", route.CategoryCreate).Methods("POST")
	s.HandleFunc("/category/{id}", route.CategoryGet).Methods("GET")
	s.HandleFunc("/category/{id}", route.CategoryUpdate).Methods("PATCH")
	s.HandleFunc("/category/{id}", route.CategoryDelete).Methods("DELETE")

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error start server: %s", err.Error()))
		return
	}
}
