package base

import (
	"errors"
	"fmt"
	"github.com/cherryReptile/Todo/internal/database"
	"github.com/cherryReptile/Todo/internal/queue"
	"github.com/cherryReptile/Todo/internal/router"
	"github.com/cherryReptile/Todo/internal/telegram"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

type App struct {
	DB               *database.SqlLite
	MuxRouter        *mux.Router
	RouterController *router.Router
	Server           *http.Server
}

func (a *App) Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("[FATAL] Not loading environment: %v", err)
	}

	db := database.Connect()
	q := queue.Run()
	service := new(telegram.Service)
	service.Init(&db)
	a.DB = &db
	a.MuxRouter = mux.NewRouter()
	a.RouterController = router.NewRouter(&q, a.DB, service)
}

func (a *App) ApiRun(port string, ch chan error) {
	a.Server = &http.Server{
		Handler:      a.MuxRouter,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	defer a.Server.Close()

	log.Printf("[DEBUG] Running server on port %s", port)

	if err := a.Server.ListenAndServe(); err != nil {
		ch <- errors.New(fmt.Sprintf("Error server: %s", err.Error()))
	}
}

func (a *App) Close() {
	if err := a.DB.DB.Close(); err != nil {
		log.Fatalf("[FATAL] Unable to close database: %v", err)
		return
	}

	if err := a.Server.Close(); err != nil {
		log.Fatalf("[FATAL] Unable to close server: %v", err)
		return
	}
}

func (a *App) POST(path string, method func(w http.ResponseWriter, r *http.Request)) {
	a.MuxRouter.HandleFunc(path, method).Methods("POST")
}

func (a *App) GET(path string, method func(w http.ResponseWriter, r *http.Request)) {
	a.MuxRouter.HandleFunc(path, method).Methods("GET")
}
