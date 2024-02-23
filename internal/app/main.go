package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"main.go/config"
	"net/http"
	"time"
)

const configPath = "config/config.json"

func startServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "LESS GOOOOOOO\n%v\n", r.Host)
		})

	server := http.Server{
		Addr:         config.Cfg.Server.Host + config.Cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  config.Cfg.Server.Timeout * time.Second,
		WriteTimeout: config.Cfg.Server.Timeout * time.Second,
	}

	logrus.Printf("starting server at %v", server.Addr)
	server.ListenAndServe()
}

func main() {
	_, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	startServer()
}
