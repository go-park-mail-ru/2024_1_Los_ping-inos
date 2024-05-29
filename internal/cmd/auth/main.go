package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"main.go/internal/auth"

	"github.com/emirpasic/gods/sets/hashset"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"main.go/config"
	Delivery "main.go/internal/auth/delivery"
	gen "main.go/internal/auth/proto"
	authRepo "main.go/internal/auth/repo"
	authUsecase "main.go/internal/auth/usecase"
	image "main.go/internal/image/protos/gen"
	. "main.go/internal/logs"
	. "main.go/internal/pkg"
)

const (
	httpPath = "../../../config/auth_http_config.yaml"
	grpcPath = "../../../config/auth_grpc_config.yaml"
)

type Delivers struct {
	http *Delivery.AuthHandler
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

	redisCli := redis.NewClient(&redis.Options{
		Addr:     httpCfg.Redis.Host + ":" + httpCfg.Redis.Port,
		Password: "",
		DB:       0,
	})

	defer func() {
		if err := redisCli.Close(); err != nil {
			logger.Logger.Errorf("Error Closing Redis connection: %v", err)
		}
		logger.Logger.Info("Redis closed without errors")
	}()

	_, pingErr := redisCli.Ping(context.Background()).Result()
	if pingErr != nil {
		logger.Logger.Errorf("Failed to ping Redis server: %v", pingErr)
	}

	grpcCfg, err := config.LoadConfig(grpcPath)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	grpcConn, err := grpc.Dial("images:50052", grpc.WithInsecure())
	if err != nil {
		logger.Logger.Fatal(err)
	}

	imageManager := image.NewImageClient(grpcConn)

	useCase := authUsecase.NewAuthUseCase(authRepo.NewAuthPersonStorage(db),
		authRepo.NewSessionStorage(redisCli), authRepo.NewInterestStorage(db), imageManager)

	srv, ok := net.Listen("tcp", grpcCfg.Server.Port)
	if ok != nil {
		logger.Logger.Fatal(err)
	}

	grpcSrever := grpc.NewServer()
	grpcDeliver := Delivery.NewGRPCDeliver(useCase)
	gen.RegisterAuthHandlServer(grpcSrever, grpcDeliver)

	httpDeliver := Delivery.NewAuthHandler(useCase)
	errors := make(chan error, 2)

	go func() {
		errors <- startServer(httpCfg, logger, Delivers{http: httpDeliver})
	}()
	go func() {
		errors <- grpcSrever.Serve(srv)
	}()

	if err = <-errors; err != nil {
		logger.Logger.Fatalf("server error: %v", err.Error())
	}
}

func startServer(cfg *config.Config, logger Log, deliver Delivers) error {
	var apiPath = cfg.ApiPath
	httpDeliver := deliver.http
	grpcConn, err := grpc.Dial("auth:50051", grpc.WithInsecure())

	if err != nil {
		return err
	}
	authManager := gen.NewAuthHandlClient(grpcConn)
	mux := http.NewServeMux()

	prometheus.MustRegister(
		auth.TotalHits,
		auth.HitDuration,
	)

	mux.Handle("/metrics", promhttp.Handler())

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
			IsAuthenticatedMiddleware(http.HandlerFunc(httpDeliver.LogoutHandler()), authManager), hashset.New("GET")),
		"logout", logger))

	mux.Handle(apiPath+"isAuth", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(httpDeliver.IsAuthenticatedHandler()), hashset.New("GET")),
		"authentication check", logger))

	mux.Handle(apiPath+"profile", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(httpDeliver.ProfileHandlers())), authManager), hashset.New("GET", "POST", "DELETE")),
		"profile", logger))

	mux.Handle(apiPath+"payment", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(httpDeliver.PaymentUrl())), authManager), hashset.New("POST")),
		"payment", logger))

	mux.Handle(apiPath+"matches", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(httpDeliver.GetMatches()), authManager), hashset.New("POST")),
		"matches", logger))

	mux.Handle(apiPath+"activatePremium", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(httpDeliver.ActivateSub()), authManager), hashset.New("POST")),
		"activate premium", logger))

	mux.Handle(apiPath+"payHistory", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(httpDeliver.GetSubHistory()), authManager), hashset.New("GET")),
		"get payment history", logger))

	metricHandler := Delivery.MetricTimeMiddleware(mux)
	server := http.Server{
		Addr:         cfg.Server.Host + cfg.Server.Port,
		Handler:      metricHandler,
		ReadTimeout:  cfg.Server.Timeout * time.Second,
		WriteTimeout: cfg.Server.Timeout * time.Second,
	}

	logger.Logger.Infof("started auth http server at %v", server.Addr)
	fmt.Printf("started auth http server at %v\n", server.Addr)
	return server.ListenAndServe()
}
