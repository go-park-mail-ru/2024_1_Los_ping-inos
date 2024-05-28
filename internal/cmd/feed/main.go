package main

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"main.go/internal/feed"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
	"main.go/config"
	gen "main.go/internal/auth/proto"
	Delivery "main.go/internal/feed/delivery"
	Repo "main.go/internal/feed/repo"
	Usecase "main.go/internal/feed/usecase"
	image "main.go/internal/image/protos/gen"
	. "main.go/internal/logs"

	"net/http"
	"time"

	. "main.go/internal/pkg"
)

const (
	httpPath = "../../../config/feed_http_config.yaml"
)

type Delivers struct {
	http *Delivery.FeedHandler
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

	grpcConn, err := grpc.Dial("images:50052", grpc.WithInsecure())
	if err != nil {
		logger.Logger.Fatal(err)
	}

	imageManager := image.NewImageClient(grpcConn)

	useCase := Usecase.New(Repo.NewPostgresStorage(db), Repo.NewWebSocStorage(), imageManager)

	grpcConn, err = grpc.Dial("auth:50051", grpc.WithInsecure())
	if err != nil {
		logger.Logger.Fatal(err)
	}
	authManager := gen.NewAuthHandlClient(grpcConn)
	httpDeliver := Delivery.NewFeedDelivery(useCase, authManager)

	errors := make(chan error, 2)
	go func() {
		errors <- startServer(httpCfg, logger, Delivers{http: httpDeliver})
	}()

	if err = <-errors; err != nil {
		logger.Logger.Fatalf("server error: %v", err.Error())
	}
}

func runSwaggerServer(cfg *config.Config, logger *Log) {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(cfg.Server.Host+cfg.Server.SwaggerPort+"/swagger/doc.json"),
	))
	err := http.ListenAndServe(config.Cfg.Server.SwaggerPort, r)
	if err != nil {
		logger.Logger.Warn(err.Error())
	}
}

func startServer(cfg *config.Config, logger Log, deliver Delivers) error {
	go runSwaggerServer(cfg, &logger)

	var apiPath = cfg.ApiPath
	feedDel := deliver.http
	authManager := feedDel.AuthManager

	mux := http.NewServeMux()

	prometheus.MustRegister(
		feed.TotalHits,
		feed.HitDuration,
	)

	mux.Handle("/metrics", promhttp.Handler())

	mux.Handle(apiPath+"cards", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(feedDel.GetCardsHandler()), authManager), hashset.New("GET")),
		"get cards", logger))
	mux.Handle(apiPath+"like", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(feedDel.CreateLike())), authManager), hashset.New("POST")),
		"like", logger))

	mux.Handle(apiPath+"dislike", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(feedDel.CreateDislike())), authManager), hashset.New("POST")),
		"dislike", logger))

	mux.Handle(apiPath+"openConnection", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(feedDel.ServeMessages()), authManager), hashset.New("GET")),
		"open connection", logger))

	mux.Handle(apiPath+"getChat", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(feedDel.GetChat()), authManager), hashset.New("POST")),
		"get chat", logger))

	mux.Handle(apiPath+"getAllChats", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(feedDel.GetAllChats()), authManager), hashset.New("GET")),
		"get all chats", logger))

	mux.Handle(apiPath+"createClaim", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(feedDel.CreateClaim()), authManager), hashset.New("POST")),
		"create claim", logger))

	mux.Handle(apiPath+"claimTypes", RequestIDMiddleware(
		AllowedMethodMiddleware(http.HandlerFunc(feedDel.GetAlClaims()), hashset.New("GET")),
		"get claim types", logger))
	metricHandler := Delivery.MetricTimeMiddleware(mux)

	server := http.Server{
		Addr:         cfg.Server.Host + cfg.Server.Port,
		Handler:      metricHandler,
		ReadTimeout:  cfg.Server.Timeout * time.Second,
		WriteTimeout: cfg.Server.Timeout * time.Second,
	}

	logger.Logger.Infof("started feed http server at %v", server.Addr)
	fmt.Printf("started feed http server at %v\n", server.Addr)
	return server.ListenAndServe()
}
