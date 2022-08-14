package main

import (
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/router"
	"log"
	"net/http"
)

func main() {
	q := queue.Run()
	sql := database.Connect()
	defer sql.DB.Close()
	route := router.NewRouter(&q, &sql)

	http.HandleFunc("/", route.Index)
	http.HandleFunc("/job", route.Test)
	http.HandleFunc("/user", route.CreateUser)
	//http.HandleFunc("/category", route.CreateCategory)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error start server: %s", err.Error()))
		return
	}
}
