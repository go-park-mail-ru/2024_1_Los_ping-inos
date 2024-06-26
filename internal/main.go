package main

import (
	"database/sql"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	. "main.go/internal/pkg"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"main.go/config"

	authDelivery "main.go/internal/auth/delivery"
	authRepo "main.go/internal/auth/repo"
	authUsecase "main.go/internal/auth/usecase"
	_ "main.go/internal/docs"
	. "main.go/internal/logs"

	profileDelivery "main.go/internal/profile/delivery"
	profileRepo "main.go/internal/profile/repo"
	profileUsecase "main.go/internal/profile/usecase"

	feedDelivery "main.go/internal/feed/delivery"
	feedRepo "main.go/internal/feed/repo"
	feedUsecase "main.go/internal/feed/usecase"

	imageDelivery "main.go/internal/image/delivery"
	imageRepo "main.go/internal/image/repo"
	imageUsecase "main.go/internal/image/usecase"
)

const configPath = "config/config.yaml"
const (
	authDeliver = iota
	profileDeliver
	feedDeliver
	imageDeliver
)

// @title SportBro API
// @version 0.1
// @host  185.241.192.216:8085
// @BasePath /api/v1/
func main() {
	logger := InitLog()

	_, err := config.LoadConfig(configPath)
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

	profilePStorage := profileRepo.NewPersonStorage(db)
	profileLStorage := profileRepo.NewLikeStorage(db)
	profileImgStorage := profileRepo.NewImageStorage(db)
	profileIntStorage := profileRepo.NewInterestStorage(db)

	feedPStorage := feedRepo.NewPersonStorage(db)
	feedLStorage := feedRepo.NewLikeStorage(db)
	feedImgStorage := feedRepo.NewImageStorage(db)
	feedIntStorage := feedRepo.NewInterestStorage(db)

	imageStorage := imageRepo.NewImageStorage(db)

	delivers := make([]interface{}, 4)
	delivers[authDeliver] = authDelivery.NewAuthHandler(authUsecase.NewAuthUseCase(authRepo.NewAuthPostgresStorage(db), authRepo.NewInterestStorage(db), authRepo.NewImageStorage(db)))
	delivers[profileDeliver] = profileDelivery.NewProfileDeliver(profileUsecase.NewProfileUseCase(profilePStorage, profileIntStorage, profileImgStorage, profileLStorage))
	delivers[feedDeliver] = feedDelivery.NewFeedDelivery(feedUsecase.New(feedPStorage, feedIntStorage, feedImgStorage, feedLStorage))
	delivers[imageDeliver] = imageDelivery.NewImageDelivery(imageUsecase.NewImageUseCase(imageStorage))

	err = StartServer(logger, delivers)
	if err != nil {
		logger.Logger.Fatalf("server error: %v", err.Error())
	}
}

func runSwaggerServer(logger *Log) {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(config.Cfg.Server.Host+config.Cfg.Server.SwaggerPort+"/swagger/doc.json"),
	))
	err := http.ListenAndServe(config.Cfg.Server.SwaggerPort, r)
	if err != nil {
		logger.Logger.Warn(err.Error())
	}
}

func StartServer(logger *Log, deliver []interface{}) error {
	go runSwaggerServer(logger)

	var apiPath = config.Cfg.ApiPath

	authResolver := deliver[authDeliver].(*authDelivery.AuthHandler).UseCase
	authDel := deliver[authDeliver].(*authDelivery.AuthHandler)
	feedDel := deliver[feedDeliver].(*feedDelivery.FeedHandler)
	imageDel := deliver[imageDeliver].(*imageDelivery.ImageHandler)
	profileDel := deliver[profileDeliver].(*profileDelivery.ProfileHandler)

	// роутер)0)
	// структура: путь, цепочка миддлвар: логирование -> методы -> [авторизация -> [CSRF]] -> функция-обработчик ручки
	mux := http.NewServeMux()
	mux.Handle(apiPath+"cards", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(feedDel.GetCardsHandler()), authResolver), hashset.New("GET")),
		"get cards", logger))

	mux.Handle(apiPath+"login", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(authDel.LoginHandler()), hashset.New("POST")),
		"login", logger))

	mux.Handle(apiPath+"registration", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(authDel.RegistrationHandler()), hashset.New("GET", "POST")),
		"registration", logger))

	mux.Handle(apiPath+"logout", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(authDel.LogoutHandler()), authResolver), hashset.New("GET")),
		"logout", logger))

	mux.Handle(apiPath+"isAuth", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(authDel.IsAuthenticatedHandler()), hashset.New("GET")),
		"authentication check", logger))

	mux.Handle(apiPath+"me", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(authDel.GetUsername()), authResolver), hashset.New("GET")),
		"username (/me)", logger))

	mux.Handle(apiPath+"getImage", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(imageDel.GetImageHandler()), authResolver), hashset.New("GET")),
		"get images", logger))

	mux.Handle(apiPath+"addImage", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(imageDel.AddImageHandler())), authResolver), hashset.New("POST")),
		"username (/me)", logger))

	mux.Handle(apiPath+"deleteImage", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(imageDel.DeleteImageHandler())), authResolver), hashset.New("POST")),
		"delete image", logger))

	mux.Handle(apiPath+"profile", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(profileDel.ProfileHandlers())), authResolver), hashset.New("GET", "POST", "DELETE")),
		"profile", logger))

	mux.Handle(apiPath+"like", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(feedDel.CreateLike())), authResolver), hashset.New("POST")),
		"like", logger))

	mux.Handle(apiPath+"matches", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(profileDel.GetMatches()), authResolver), hashset.New("GET")),
		"matches", logger))

	mux.Handle(apiPath+"dislike", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(feedDel.CreateDislike())), authResolver), hashset.New("POST")),
		"dislike", logger))

	server := http.Server{
		Addr:         config.Cfg.Server.Host + config.Cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  config.Cfg.Server.Timeout * time.Second,
		WriteTimeout: config.Cfg.Server.Timeout * time.Second,
	}

	logger.Logger.Infof("started server at %v", server.Addr)
	fmt.Printf("started server at %v\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
