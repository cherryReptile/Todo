package main

import (
	"github.com/cherryReptile/Todo/internal/base"
	"log"
	"os"
)

func main() {
	errChan := make(chan error, 1)

	app := new(base.App)
	app.Init()
	app.POST("/", app.RouterController.Start)

	go app.ApiRun("3000", errChan)

	err := <-errChan
	if err != nil {
		app.Close()
		log.Printf("[FATAL] %v", err)
		os.Exit(1)
	}
}
