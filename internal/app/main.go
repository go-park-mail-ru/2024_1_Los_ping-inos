package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	psqInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.Get("database.host"), viper.Get("database.port"), viper.Get("database.user"),
		viper.Get("database.password"), viper.Get("database.dbname"))

	db, err := sql.Open("postgres", psqInfo)
	if err != nil {
		logrus.Fatalf("can't open db! %v", err.Error())
	}
	if err := db.Ping(); err != nil {
		println(err.Error())
		logrus.Fatal(err)
	}
	defer db.Close()

	personStore := storage.NewPersonStorage(db)
	auth := service.NewAuthHandler(personStore)
	serv := service.New(personStore)
	deliver := delivery.New(serv, auth)
	err = delivery.StartServer(deliver)
	if err != nil {
		logrus.Fatal(err)
	}
}