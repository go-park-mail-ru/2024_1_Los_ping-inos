package delivery

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"main.go/config"
	. "main.go/internal/logs"
	"net/http"
	"sync/atomic"
	"time"
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

func runSwaggerServer() {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(config.Cfg.Server.Host+config.Cfg.Server.SwaggerPort+"/swagger/doc.json"),
	))
	err := http.ListenAndServe(config.Cfg.Server.SwaggerPort, r)
	if err != nil {
		Log.Warn(err.Error())
	}
}

func StartServer(deliver ...*Deliver) error {
	go runSwaggerServer()

	var apiPath = config.Cfg.ApiPath

	// роутер)0)
	// структура: путь, цепочка миддлвар: авторизация(методы(функция-обработчик ручки))
	mux := http.NewServeMux()
	mux.Handle(apiPath+"cards", AllowedMethodMiddleware(
		IsAuthenticatedMiddleware(http.HandlerFunc(deliver[0].GetCardsHandler()), deliver[0]),
		hashset.New("GET")))
	mux.Handle(apiPath+"login", AllowedMethodMiddleware(
		http.HandlerFunc(deliver[0].LoginHandler()),
		hashset.New("POST")))
	mux.Handle(apiPath+"registration", AllowedMethodMiddleware(
		http.HandlerFunc(deliver[0].RegistrationHandler()),
		hashset.New("GET", "POST")))
	mux.Handle(apiPath+"logout", AllowedMethodMiddleware(
		IsAuthenticatedMiddleware(http.HandlerFunc(deliver[0].LogoutHandler()), deliver[0]),
		hashset.New("GET")))
	mux.Handle(apiPath+"isAuth", AllowedMethodMiddleware(
		http.HandlerFunc(deliver[0].IsAuthenticatedHandler()), hashset.New("GET")))
	mux.Handle(apiPath+"me", AllowedMethodMiddleware(
		IsAuthenticatedMiddleware(http.HandlerFunc(deliver[0].GetUsername()), deliver[0]),
		hashset.New("GET")))
	mux.Handle(apiPath+"profile", AllowedMethodMiddleware(
		IsAuthenticatedMiddleware(http.HandlerFunc(deliver[0].ProfileHandlers()), deliver[0]),
		hashset.New("GET", "POST", "DELETE")))

	server := http.Server{
		Addr:         config.Cfg.Server.Host + config.Cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  config.Cfg.Server.Timeout * time.Second,
		WriteTimeout: config.Cfg.Server.Timeout * time.Second,
	}

	Log.Infof("started server at %v", server.Addr)
	fmt.Printf("started server at %v\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
