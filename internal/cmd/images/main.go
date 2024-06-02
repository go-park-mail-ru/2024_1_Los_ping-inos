package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"main.go/config"
	grpcDelivery "main.go/internal/image/delivery/grpc"
	httpDelivery "main.go/internal/image/delivery/http"
	gen "main.go/internal/image/protos/gen"
	"main.go/internal/image/usecase"
	. "main.go/internal/logs"
)

const (
	httpPath = "../../../config/image_http_config.yaml"
	grpcPath = "../../../config/image_grpc_config.yaml"
)

func main() {
	logger := InitLog()

	_, err := config.LoadConfig(httpPath)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	psqInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.Get("database.host"), viper.Get("database.port"), viper.Get("database.user"),
		viper.Get("database.password"), viper.Get("database.dbname"))
	//config, err := config.ReadConfig()
	//println("CONFIG", config.Database)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
		return
	}

	core, err := usecase.GetCore(psqInfo)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
		return
	}

	api := httpDelivery.GetApi(core, logger)

	grpcCfg, err := config.LoadConfig(grpcPath)
	if err != nil {
		logger.Logger.Fatal(err)
		println("princess")
	}
	println("THIS IS GRPC PORT", grpcCfg.Server.Port)
	srv, ok := net.Listen("tcp", grpcCfg.Server.Port)
	if ok != nil {
		logger.Logger.Fatal(err)
		println("princess1")
	}
	grpcServer := grpc.NewServer() // TODO интерсептеры для метрик сюда
	grpcDeliver := grpcDelivery.NewGRPCDeliver(core)
	gen.RegisterImageServer(grpcServer, grpcDeliver)

	errs := make(chan error, 2)

	go func() {
		errs <- api.ListenAndServe()
	}()
	go func() {
		errs <- grpcServer.Serve(srv)
	}()

	err = <-errs
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
	}
}
