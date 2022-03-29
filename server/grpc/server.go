package main

import (
	"fmt"
	"github.com/joho/godotenv"
	grpc2 "github.com/norby7/shortening-service/interfaceAdapters/grpc"
	"github.com/norby7/shortening-service/interfaceAdapters/grpc/protocol"
	"github.com/norby7/shortening-service/usecases/repository"
	ucCache "github.com/norby7/shortening-service/usecases/repository/cache"
	"github.com/norby7/shortening-service/usecases/repository/storage"
	ucService "github.com/norby7/shortening-service/usecases/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// StartServer starts a new grpc server and registers the UrlServiceServer to it
func StartServer(port *int, service ucService.Interactor, logger *log.Logger) {
	urlService := grpc2.NewUrlGrpcService(service, logger)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	//  register  grpcurl  The required  reflection  service
	reflection.Register(grpcServer)

	protocol.RegisterUrlServiceServer(grpcServer, urlService)
	go func() {
		fmt.Printf("Starting grpc server on port: %d\n", *port)
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Received terminate, graceful shutdown", sig)

	grpcServer.GracefulStop()
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

	port := os.Getenv("GRPC_PORT")

	portAdr, err := strconv.Atoi(port)
	if err != nil {
		portAdr = 3000
	}

	StartServer(&portAdr, service, l)
}
