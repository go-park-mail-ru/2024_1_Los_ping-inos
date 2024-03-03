package delivery

import (
	models "main.go/db"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"main.go/config"
)

type Service interface {
	GetCards(sessionID string) ([]models.Person, error)
}

type Auth interface {
	IsAuthenticated(sessionID string) bool
	Login(email, password string) (string, error)
	Logout(sessionID string) error
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
