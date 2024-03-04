package delivery

import (
	"github.com/sirupsen/logrus"
	"main.go/config"
	"main.go/internal/types"
	"net/http"
	"time"
)

type Service interface {
	GetCards(sessionID string, firstID types.UserID) (string, error)
}

type Auth interface {
	IsAuthenticated(sessionID string) bool
	Login(email, password string) (string, error)
	Logout(sessionID string) error
	Registration(Name string, Birthday string, Gender string, Email string, Password string) error
}

type Deliver struct {
	serv Service
	auth Auth
}

func New(service Service, auth Auth) *Deliver {
	return &Deliver{serv: service, auth: auth}
}

// StartServer - запуск сервера
func StartServer(deliver ...*Deliver) error {
	mux := http.NewServeMux()

	// тут хендлеры добавлять
	deliver[0].GetCardsHandler(mux)
	deliver[0].GetLoginHandler(mux)
	deliver[0].GetRegistrationHandler(mux)
	deliver[0].GetLogoutHandler(mux)

	server := http.Server{
		Addr:         config.Cfg.Server.Host + config.Cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  config.Cfg.Server.Timeout * time.Second,
		WriteTimeout: config.Cfg.Server.Timeout * time.Second,
	}

	logrus.Printf("starting server at %v", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
