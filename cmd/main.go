package main

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/router"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	godotenv.Load(".env")
	q := queue.Run()
	sql := database.Connect()
	defer sql.DB.Close()
	route := router.NewRouter(&q, &sql)
	r := mux.NewRouter()
	s := r.Host("127.0.0.1:3000").Subrouter()

	s.HandleFunc("/", route.Index).Methods("GET")
	s.HandleFunc("/test", route.Start)
	s.HandleFunc("/category", route.CategoryCreate)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error start server: %s", err.Error()))
		return
	}
}
