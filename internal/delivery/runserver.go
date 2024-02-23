package delivery

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"main.go/config"
)

type Service interface {
	GetCoolIdsList() ([]string, error)
}

type Deliver struct {
	serv Service
}

func New(l Service) *Deliver {
	return &Deliver{serv: l}
}

// StartServer - запуск сервера
func StartServer(deliver *Deliver) error {
	mux := http.NewServeMux()

	// тут хендлеры добавлять
	deliver.landingHandler(mux)

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
