package delivery

import (
	"github.com/sirupsen/logrus"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"net/http"
)

func (deliver *Deliver) CreateDislike() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		Log.WithFields(logrus.Fields{RequestID: request.Context().Value(RequestID).(int64)}).Warn("not implemented dislike")
		requests.SendResponse(respWriter, request, http.StatusOK, nil)
	}
}
