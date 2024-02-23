package main

import (
	"github.com/sirupsen/logrus"
	"main.go/config"
	"main.go/internal/delivery"
)

const configPath = "config/config.yaml"

func main() {
	_, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	err = delivery.StartServer()
	if err != nil {
		logrus.Fatal(err)
	}
}
