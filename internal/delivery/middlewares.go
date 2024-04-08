package delivery

import (
	"context"
	"net/http"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/sirupsen/logrus"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
)

func IsAuthenticatedMiddleware(next http.Handler, deliver *Deliver) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		session, err := request.Cookie("session_id")                                         // проверка авторизации
		if err != nil || session == nil || !deliver.auth.IsAuthenticated(session.Value, 0) { // я хз, как сделать один id у мидлвары и следующего обработчика
			requests.SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(respWriter, request)
	})
}

func AllowedMethodMiddleware(next http.Handler, methods *hashset.Set) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodOptions {
			requests.SendResponse(respWriter, request, http.StatusOK, nil)
			return
		}

		if !methods.Contains(request.Method) {
			requests.SendResponse(respWriter, request, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		next.ServeHTTP(respWriter, request)
	})
}

func RequestIDMiddleware(next http.Handler, deliver *Deliver, msg string) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		requestID := deliver.nextRequest()
		contexted := request.WithContext(context.WithValue(context.Background(), RequestID, requestID))
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info(msg + " request")
		next.ServeHTTP(respWriter, contexted)
	})
}
