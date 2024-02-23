package delivery

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"main.go/config"
)

func StartServer() error {
	mux := http.NewServeMux()

	// сюда добавлять хендлеры страничек
	mux.HandleFunc("/", landing)

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
