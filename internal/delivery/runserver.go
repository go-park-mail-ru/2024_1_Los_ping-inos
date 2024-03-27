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
	// "сырой" mux
	rawMux := http.NewServeMux()
	rawMux.HandleFunc(apiPath+"cards", deliver[0].GetCardsHandler())
	rawMux.HandleFunc(apiPath+"login", deliver[0].LoginHandler())
	rawMux.HandleFunc(apiPath+"registration", deliver[0].RegistrationHandler())
	rawMux.HandleFunc(apiPath+"logout", deliver[0].LogoutHandler())
	rawMux.HandleFunc(apiPath+"isAuth", deliver[0].IsAuthenticatedHandler())
	rawMux.HandleFunc(apiPath+"me", deliver[0].GetUsername())
	rawMux.HandleFunc(apiPath+"profile", deliver[0].UpdateProfile())

	// обёртки миддлвар на методы и авторизованность
	authHandler := IsAuthenticatedMiddleware(rawMux, deliver[0])

	cardsHandler := AllowedMethodMiddleware(authHandler, hashset.New("GET"))
	loginHandler := AllowedMethodMiddleware(rawMux, hashset.New("POST"))
	registrationHandler := AllowedMethodMiddleware(rawMux, hashset.New("GET", "POST"))
	logoutHandler := AllowedMethodMiddleware(authHandler, hashset.New("GET"))
	isAuthHandler := AllowedMethodMiddleware(rawMux, hashset.New("GET"))
	usernameHandler := AllowedMethodMiddleware(authHandler, hashset.New("GET"))
	profileUpdateHandler := AllowedMethodMiddleware(authHandler, hashset.New("POST"))

	// сохранение обёрток
	mux := http.NewServeMux()
	mux.Handle(apiPath+"cards", cardsHandler)
	mux.Handle(apiPath+"login", loginHandler)
	mux.Handle(apiPath+"registration", registrationHandler)
	mux.Handle(apiPath+"logout", logoutHandler)
	mux.Handle(apiPath+"isAuth", isAuthHandler)
	mux.Handle(apiPath+"me", usernameHandler)
	mux.Handle(apiPath+"profile", profileUpdateHandler)

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
