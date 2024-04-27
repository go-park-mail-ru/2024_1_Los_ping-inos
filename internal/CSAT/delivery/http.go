package delivery

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"main.go/internal/CSAT"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
)

type HttpHandler struct {
	useCase CSAT.UseCase
}

func NewHttpHandler(uc CSAT.UseCase) *HttpHandler {
	return &HttpHandler{
		useCase: uc,
	}
}

func (deliver *HttpHandler) Create() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)

		body, err := io.ReadAll(request.Body)
		if err != nil { // TODO эти два блока вынести в отдельную функцию и напихать её во все ручки
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}
		var requestBody CSAT.CreateRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		err = deliver.useCase.Create(request.Context(), requestBody)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("can't create rate: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent created rate")
		requests.SendResponse(respWriter, request, http.StatusOK, nil)
	}
}

func (deliver *HttpHandler) Read() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)

		//body, err := io.ReadAll(request.Body)
		// if err != nil { // TODO эти два блока вынести в отдельную функцию и напихать её во все ручки
		// 	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
		// 	requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		// 	return
		// }

		// var requestBody CSAT.StatRequest
		// err = json.Unmarshal(body, &requestBody)
		// if err != nil {
		// 	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
		// 	requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		// 	return
		// }

		allQ, err := deliver.useCase.GetStat(request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("can't create rate: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent created rate")
		requests.SendResponse(respWriter, request, http.StatusOK, allQ)
	}
}
