package main

import (
	"database/sql"
	"fmt"
	"log"

	//"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"main.go/config"
	_ "main.go/db"
	"main.go/internal/delivery"
	_ "main.go/internal/docs"
	. "main.go/internal/logs"
	"main.go/internal/service"
	"main.go/internal/storage"
)

const configPath = "../config/config.yaml"

const (
	vkCloudHotboxEndpoint = "https://hb.vkcs.cloud"
	defaultRegion         = "ru-msk"
)

// @title SportBro API
// @version 0.1
// @host  185.241.192.216:8081
// @BasePath /
func main() {
	InitLog()

	_, err := config.LoadConfig(configPath)
	if err != nil {
		Log.Fatal(err)
	}
	psqInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.Get("database.host"), viper.Get("database.port"), viper.Get("database.user"),
		viper.Get("database.password"), viper.Get("database.dbname"))

	db, err := sql.Open("postgres", psqInfo)
	if err != nil {
		Log.Fatalf("can't open db: %v", err.Error())
	}
	if err = db.Ping(); err != nil {
		println(err.Error())
		Log.Fatal(err)
	}
	defer db.Close()

	sess, _ := session.NewSession()
	svc := s3.New(sess, aws.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))

	if res, err := svc.ListBuckets(nil); err != nil {
		log.Fatalf("Unable to list buckets, %v", err)
	} else {
		for _, b := range res.Buckets {
			log.Printf("* %s created on %s \n", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
		}
	}

	personStore := storage.NewPersonStorage(db)
	interestStore := storage.NewInterestStorage(db)
	imageStore := storage.NewImageStorage(db)

	auth := service.NewAuthHandler(personStore)
	serv := service.New(personStore, interestStore, imageStore)
	deliver := delivery.New(serv, auth)
	err = delivery.StartServer(deliver)
	if err != nil {
		Log.Fatalf("server error: %v", err.Error())
	}
}
