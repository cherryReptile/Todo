package main

import (
	"fmt"
	"github.com/pavel-one/GoStarter/internal/database"
	"github.com/pavel-one/GoStarter/internal/queue"
	"github.com/pavel-one/GoStarter/internal/router"
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
	http.HandleFunc("/create", route.Create)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error start server: %s", err.Error()))
		return
	}
}
