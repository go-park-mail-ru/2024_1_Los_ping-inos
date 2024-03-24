package delivery

import (
	"fmt"
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

	// "сырой" mux
	rawMux := http.NewServeMux()
	rawMux.HandleFunc("/cards", deliver[0].GetCardsHandler())
	rawMux.HandleFunc("/login", deliver[0].LoginHandler())
	rawMux.HandleFunc("/registration", deliver[0].RegistrationHandler())
	rawMux.HandleFunc("/logout", deliver[0].LogoutHandler())
	rawMux.HandleFunc("/isAuth", deliver[0].IsAuthenticatedHandler())
	rawMux.HandleFunc("/me", deliver[0].GetUsername())

	// обёртки миддлвар на методы и авторизованность
	authHandler := IsAuthenticatedMiddleware(rawMux, deliver[0])

	cardsHandler := AllowedMethodMiddleware(authHandler, map[string]struct{}{"GET": {}})
	loginHandler := AllowedMethodMiddleware(rawMux, map[string]struct{}{"POST": {}})
	registrationHandler := AllowedMethodMiddleware(rawMux, map[string]struct{}{"GET": {}, "POST": {}})
	logoutHandler := AllowedMethodMiddleware(authHandler, map[string]struct{}{"GET": {}})
	isAuthHandler := AllowedMethodMiddleware(rawMux, map[string]struct{}{"GET": {}})
	usernameHandler := AllowedMethodMiddleware(authHandler, map[string]struct{}{"GET": {}})

	// сохранение обёрток
	mux := http.NewServeMux()
	mux.Handle("/cards", cardsHandler)
	mux.Handle("/login", loginHandler)
	mux.Handle("/registration", registrationHandler)
	mux.Handle("/logout", logoutHandler)
	mux.Handle("/isAuth", isAuthHandler)
	mux.Handle("/me", usernameHandler)

	server := http.Server{
		Addr:         config.Cfg.Server.Host + config.Cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  config.Cfg.Server.Timeout * time.Second,
		WriteTimeout: config.Cfg.Server.Timeout * time.Second,
	}

	Log.Infof("started server at %v\n", server.Addr)
	fmt.Printf("started server at %v\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
