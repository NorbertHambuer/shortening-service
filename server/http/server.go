package main

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpC "github.com/norby7/shortening-service/interfaceAdapters/http"
	"github.com/norby7/shortening-service/usecases/repository"
	ucCache "github.com/norby7/shortening-service/usecases/repository/cache"
	"github.com/norby7/shortening-service/usecases/repository/storage"
	ucService "github.com/norby7/shortening-service/usecases/service"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// RegisterRoutes registers the http server routes
func RegisterRoutes(r *mux.Router, c httpC.Controller) {
	r.HandleFunc("/api", c.Add).Methods("POST")
	r.HandleFunc("/api/{code:[a-zA-Z0-9]+}", c.Delete).Methods("DELETE")
	r.HandleFunc("/api/{code:[a-zA-Z0-9]+}", c.Get).Methods("GET")

	// create Redoc configuration
	ops := middleware.RedocOpts{
		SpecURL: "/swagger.yaml",
	}

	// add swagger documentation routes
	sh := middleware.Redoc(ops, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	r.HandleFunc("/counter/{code:[a-zA-Z0-9]+}", c.GetCounter).Methods("GET")
	r.HandleFunc("/{code:[a-zA-Z0-9]+}", c.RedirectShortUrl).Methods("GET")
}

// StartServer starts a new http server that listens on the given port
func StartServer(r *mux.Router, port int) {
	s := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// start server on a different goroutine
	go func() {
		log.Println("Starting server on port " + strconv.Itoa(port))

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(fmt.Sprintf("unable to start http server: %s", err.Error()))
		}

	}()

	// create a signal channel that will be notified for Interrupt and Kill signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// wait for a signal
	sig := <-sigChan
	log.Println("Received terminate, graceful shutdown", sig)

	// create context with timeout, the server will wait 30 seconds for all connections to finish
	tc, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := s.Shutdown(tc); err != nil {
		log.Fatalln(fmt.Sprintf("error shuting down server: %s", err.Error()))
	}

}

func main() {
	// set the random seed
	rand.Seed(time.Now().UnixNano())

	// load .env file
	godotenv.Load()
	workersStr := os.Getenv("COUNTER_WORKERS")

	workers, err := strconv.Atoi(workersStr)
	if err != nil {
		workers = 10
	}

	dbPath := "./database/sqlite/urls.db"
	l := log.New(os.Stdout, "urls-api", log.LstdFlags)

	// create sqlite database file if it doesn't exists
	err = storage.CreateDatabase(dbPath)
	if err != nil {
		l.Fatalln(err.Error())
	}

	// new sqlite repository
	sqliteStorage, err := storage.NewSqliteStorage(dbPath, workers)
	if err != nil {
		l.Fatalln("unable to create new repository: " + err.Error())
	}

	defer sqliteStorage.Handler.Close()

	// checks if the schema exists and initialize it if it doesn't
	err = storage.ValidateSchema(sqliteStorage.Handler)
	if err != nil {
		l.Fatalln(err.Error())
	}

	// creates a new cache object
	redisCache, err := ucCache.NewRedisCache(os.Getenv("REDIS_HOSTNAME"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))
	if err != nil {
		l.Println("unable to connect to redis cache: " + err.Error())
	}

	urlRepo := repository.NewUrlRepository(sqliteStorage, redisCache, l)

	service := ucService.NewService(urlRepo, workers, os.Getenv("REDIRECT_DOMAIN"))
	controller := httpC.NewController(service, l)

	muxRouter := mux.NewRouter()
	RegisterRoutes(muxRouter, *controller)

	port := os.Getenv("PORT")

	portAdr, err := strconv.Atoi(port)
	if err != nil {
		portAdr = 3000
	}

	StartServer(muxRouter, portAdr)
}
