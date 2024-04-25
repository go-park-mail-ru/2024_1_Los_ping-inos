package main

import (
	"database/sql"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"main.go/config"
	httpDeliverivery "main.go/internal/auth/delivery"
	authRepo "main.go/internal/auth/repo"
	authUsecase "main.go/internal/auth/usecase"
	. "main.go/internal/logs"
	. "main.go/internal/pkg"
	"net/http"
	"time"
)

const configPath = "config/auth_http_config.yaml"

type Delivers struct {
	http *httpDeliverivery.AuthHandler
}

func main() {
	logger := InitLog()

	httpCfg, err := config.LoadConfig(configPath)
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

	httpDeliver := httpDeliverivery.NewAuthHandler(authUsecase.NewAuthUseCase(authRepo.NewAuthPersonStorage(db), authRepo.NewInterestStorage(db), authRepo.NewImageStorage(db)))

	errors := make(chan error, 2)
	go func() {
		errors <- startServer(httpCfg, logger, Delivers{http: httpDeliver})
	}()

	if err = <-errors; err != nil {
		logger.Logger.Fatalf("server error: %v", err.Error())
	}
}

func startServer(cfg *config.Config, logger *Log, deliver Delivers) error {
	var apiPath = config.Cfg.ApiPath
	httpDeliver := deliver.http
	authResolver := httpDeliver.UseCase
	mux := http.NewServeMux()
	mux.Handle(apiPath+"login", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(httpDeliver.LoginHandler()), hashset.New("POST")),
		"login", logger))

	mux.Handle(apiPath+"registration", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(httpDeliver.RegistrationHandler()), hashset.New("GET", "POST")),
		"registration", logger))

	mux.Handle(apiPath+"logout", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(httpDeliver.LogoutHandler()), authResolver), hashset.New("GET")),
		"logout", logger))

	mux.Handle(apiPath+"isAuth", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(httpDeliver.IsAuthenticatedHandler()), hashset.New("GET")),
		"authentication check", logger))

	mux.Handle(apiPath+"me", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(httpDeliver.GetUsername()), authResolver), hashset.New("GET")),
		"username (/me)", logger))

	mux.Handle(apiPath+"profile", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(httpDeliver.ProfileHandlers())), authResolver), hashset.New("GET", "POST", "DELETE")),
		"profile", logger))

	mux.Handle(apiPath+"matches", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(httpDeliver.GetMatches()), authResolver), hashset.New("GET")),
		"matches", logger))

	server := http.Server{
		Addr:         cfg.Server.Host + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.Timeout * time.Second,
		WriteTimeout: cfg.Server.Timeout * time.Second,
	}

	logger.Logger.Infof("started auth http server at %v", server.Addr)
	fmt.Printf("started auth http server at %v\n", server.Addr)
	err := server.ListenAndServe()

	return err
}
