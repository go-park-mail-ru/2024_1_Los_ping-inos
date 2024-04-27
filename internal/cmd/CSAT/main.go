package main

import (
	"database/sql"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"main.go/config"
	Delivery "main.go/internal/CSAT/delivery"
	Repo "main.go/internal/CSAT/repo"
	Usecase "main.go/internal/CSAT/usecase"
	gen "main.go/internal/auth/proto"
	. "main.go/internal/logs"
	. "main.go/internal/pkg"
	"net/http"
	"time"
)

const (
	httpPath = "config/csat_http_config.yaml"
)

type Delivers struct {
	http *Delivery.HttpHandler
}

func main() {
	logger := InitLog()

	httpCfg, err := config.LoadConfig(httpPath)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	psqInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.Get("database.host"), viper.Get("database.port"), viper.Get("database.user"),
		viper.Get("database.password"), viper.Get("database.dbname"))

	db, err := sql.Open("postgres", psqInfo)
	if err != nil {
		logger.Logger.Fatalf("can't open db: %v", err.Error())
	}
	if err = db.Ping(); err != nil {
		println(err.Error())
		logger.Logger.Fatal(err)
	}
	defer db.Close()

	useCase := Usecase.NewUseCase(Repo.NewCsatStorage(db))

	httpDeliver := Delivery.NewHttpHandler(useCase)

	err = startServer(httpCfg, logger, Delivers{http: httpDeliver})
	logger.Logger.Fatalf("server error: %v", err.Error())
}

func startServer(cfg *config.Config, logger Log, deliver Delivers) error {
	var apiPath = cfg.ApiPath
	httpDeliver := deliver.http

	grpcConn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		return err
	}
	authManager := gen.NewAuthHandlClient(grpcConn)
	mux := http.NewServeMux()

	mux.Handle(apiPath+"createRate", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(httpDeliver.Create()), authManager), hashset.New("POST")),
		"create rate", logger))

	//mux.Handle(apiPath+"isAuth", RequestIDMiddleware(
	//	AllowedMethodMiddleware(
	//		http.HandlerFunc(httpDeliver.IsAuthenticatedHandler()), hashset.New("GET")),
	//	"authentication check", logger))

	server := http.Server{
		Addr:         cfg.Server.Host + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.Timeout * time.Second,
		WriteTimeout: cfg.Server.Timeout * time.Second,
	}

	logger.Logger.Infof("started CSAT http server at %v", server.Addr)
	fmt.Printf("started CSAT http server at %v\n", server.Addr)
	return server.ListenAndServe()
}
