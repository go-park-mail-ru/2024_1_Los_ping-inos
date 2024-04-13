package delivery

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"main.go/config"
	. "main.go/internal/logs"
)

type Deliver struct {
	serv          Service
	auth          Auth
	lastRequestID int64
}

func New(service Service, auth Auth) *Deliver {
	return &Deliver{serv: service, auth: auth}
}

func (deliver *Deliver) nextRequest() int64 {
	atomic.AddInt64(&deliver.lastRequestID, 1)
	return deliver.lastRequestID
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

func StartServer(logger *Log, deliver ...*Deliver) error {
	go runSwaggerServer(logger)

	var apiPath = config.Cfg.ApiPath

	// роутер)0)
	// структура: путь, цепочка миддлвар: логирование -> методы -> [авторизация -> [CSRF]] -> функция-обработчик ручки
	mux := http.NewServeMux()
	mux.Handle(apiPath+"cards", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(deliver[0].GetCardsHandler()), deliver[0]), hashset.New("GET")),
		deliver[0], "get cards", logger))

	mux.Handle(apiPath+"login", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(deliver[0].LoginHandler()), hashset.New("POST")),
		deliver[0], "login", logger))

	mux.Handle(apiPath+"registration", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(deliver[0].RegistrationHandler()), hashset.New("GET", "POST")),
		deliver[0], "registration", logger))

	mux.Handle(apiPath+"logout", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(deliver[0].LogoutHandler()), deliver[0]), hashset.New("GET")),
		deliver[0], "logout", logger))

	mux.Handle(apiPath+"isAuth", RequestIDMiddleware(
		AllowedMethodMiddleware(
			http.HandlerFunc(deliver[0].IsAuthenticatedHandler()), hashset.New("GET")),
		deliver[0], "authentication check", logger))

	mux.Handle(apiPath+"me", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(deliver[0].GetUsername()), deliver[0]), hashset.New("GET")),
		deliver[0], "username (/me)", logger))

	mux.Handle(apiPath+"getImage", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(deliver[0].GetImageHandler()), deliver[0]), hashset.New("GET")),
		deliver[0], "get images", logger))

	mux.Handle(apiPath+"addImage", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(deliver[0].AddImageHandler())), deliver[0]), hashset.New("POST")),
		deliver[0], "username (/me)", logger))

	mux.Handle(apiPath+"deleteImage", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(deliver[0].DeleteImageHandler())), deliver[0]), hashset.New("POST")),
		deliver[0], "delete image", logger))

	mux.Handle(apiPath+"profile", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(deliver[0].ProfileHandlers())), deliver[0]), hashset.New("GET", "POST", "DELETE")),
		deliver[0], "profile", logger))

	mux.Handle(apiPath+"like", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(deliver[0].CreateLike())), deliver[0]), hashset.New("POST")),
		deliver[0], "like", logger))

	mux.Handle(apiPath+"matches", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(http.HandlerFunc(deliver[0].GetMatches()), deliver[0]), hashset.New("GET")),
		deliver[0], "matches", logger))

	mux.Handle(apiPath+"dislike", RequestIDMiddleware(
		AllowedMethodMiddleware(
			IsAuthenticatedMiddleware(CSRFMiddleware(http.HandlerFunc(deliver[0].CreateDislike())), deliver[0]), hashset.New("POST")),
		deliver[0], "dislike", logger))

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
