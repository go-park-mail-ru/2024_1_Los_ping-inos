package main

import (
	"github.com/sirupsen/logrus"
	"main.go/config"
	delivery "main.go/internal/image/delivery/http"
	"main.go/internal/image/usecase"
	. "main.go/internal/logs"
)

// const (
// 	httpPath = "config/image_http_config.yaml"
// 	grpcPath = "config/auth_grpc_config.yaml"
// )

const (
	httpPath = "../../../config/image_http_config.yaml"
	grpcPath = "../../../config/image_grpc_config.yaml"
)

func main() {
	logger := InitLog()

	httpCfg, err := config.LoadConfig(httpPath)

	//config, err := config.ReadConfig()
	//println("CONFIG", config.Database)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
		return
	}
	//var ctx context.Context

	core, err := usecase.GetCore(config)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
		return
	}

	api := delivery.GetApi(core, logger)

	errs := make(chan error, 2)

	go func() {
		errs <- api.ListenAndServe()
	}()

	err = <-errs
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
	}
}
