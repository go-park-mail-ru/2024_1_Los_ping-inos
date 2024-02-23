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

	stor := storage.Storage{}
	serv := service.New(stor)
	deliver := delivery.New(serv)

	err = delivery.StartServer(deliver)
	if err != nil {
		logrus.Fatal(err)
	}
}
