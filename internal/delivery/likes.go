package delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	. "main.go/config"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	"net/http"
)

func (deliver *Deliver) CreateLike() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		requestID := request.Context().Value(RequestID).(int64)

		body, err := io.ReadAll(request.Body)
		if err != nil { // TODO эти два блока вынести в отдельную функцию и напихать её во все ручки
			Log.WithFields(logrus.Fields{RequestID: requestID}).Info("bad body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		var requestBody requests.CreateLikeRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		err = deliver.serv.CreateLike(request.Context().Value(RequestUserID).(types.UserID), requestBody.Profile2, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't update profile: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, nil)
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("update profile sent response")
	}
}
