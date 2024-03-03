package main

import (
	"github.com/sirupsen/logrus"
	"main.go/config"
	"main.go/internal/delivery"
	"main.go/internal/service"
	"main.go/internal/storage"
)

const configPath = "config/config.yaml"

func main() {
	_, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	personStore := storage.PersonStorage{}
	auth := service.NewAuthHandler()
	serv := service.New(&personStore)
	deliver := delivery.New(serv, auth)

	err = delivery.StartServer(deliver)
	if err != nil {
		logrus.Fatal(err)
	}
}
