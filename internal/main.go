package main

import (
	"database/sql"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"google.golang.org/grpc"
	. "main.go/internal/pkg"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"main.go/config"

	authGrpc "main.go/internal/auth/proto"
	_ "main.go/internal/docs"
	. "main.go/internal/logs"

	imageDelivery "main.go/internal/image/delivery"
	imageRepo "main.go/internal/image/repo"
	imageUsecase "main.go/internal/image/usecase"
)

const configPath = "config/config.yaml"
const (
	imageDeliver = iota
)

// @title SportBro API
// @version 0.1
// @host  185.241.192.216:8085
// @BasePath /api/v1/
func main() {
	logger := InitLog()

	cfg, err := config.LoadConfig(configPath)
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

	imageStorage := imageRepo.NewImageStorage(db)

	delivers := make([]interface{}, 4)
	delivers[imageDeliver] = imageDelivery.NewImageDelivery(imageUsecase.NewImageUseCase(imageStorage))

	err = StartServer(cfg, logger, delivers)
	if err != nil {
		logger.Logger.Fatalf("server error: %v", err.Error())
	}
}

func StartServer(cfg *config.Config, logger Log, deliver []interface{}) error {

	var apiPath = cfg.ApiPath
	imageDel := deliver[imageDeliver].(*imageDelivery.ImageHandler)

	grpcConn, err := grpc.Dial("auth:50051", grpc.WithInsecure())
	if err != nil {
		return err
	}

	authManager := authGrpc.NewAuthHandlClient(grpcConn)
	// роутер)0)
	// структура: путь, цепочка миддлвар: логирование -> методы -> [авторизация -> [CSRF]] -> функция-обработчик ручки
	mux := http.NewServeMux()

	mux.Handle(apiPath+"getImage", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(imageDel.GetImageHandler()), authManager), hashset.New("GET")),
		"get images", logger))

	mux.Handle(apiPath+"addImage", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(imageDel.AddImageHandler())), authManager), hashset.New("POST")),
		"username (/me)", logger))

	mux.Handle(apiPath+"deleteImage", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(imageDel.DeleteImageHandler())), authManager), hashset.New("POST")),
		"delete image", logger))

	server := http.Server{
		Addr:         config.Cfg.Server.Host + config.Cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  config.Cfg.Server.Timeout * time.Second,
		WriteTimeout: config.Cfg.Server.Timeout * time.Second,
	}

	logger.Logger.Infof("started server at %v", server.Addr)
	fmt.Printf("started server at %v\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
