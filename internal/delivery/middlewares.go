package delivery

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	"net/http"
)

const CSRFHeader = "csrft"

func IsAuthenticatedMiddleware(next http.Handler, deliver *Deliver) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		log := request.Context().Value(Logg).(Log)

		session, err := request.Cookie("session_id") // проверка авторизации
		if err != nil || session == nil {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("unauthorized")
			requests.SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized")
			return
		}
		id, authorized := deliver.auth.IsAuthenticated(session.Value, request.Context())

		if !authorized {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("unauthorized")
			requests.SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized")
			return
		}
		log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("authorized")
		contexted := request.WithContext(context.WithValue(request.Context(), RequestUserID, id))
		contexted = request.WithContext(context.WithValue(contexted.Context(), RequestSID, session.Value))
		next.ServeHTTP(respWriter, contexted)
	})
}

func AllowedMethodMiddleware(next http.Handler, methods *hashset.Set) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		log := request.Context().Value(Logg).(Log)

		if request.Method == http.MethodOptions {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("preflight")
			requests.SendResponse(respWriter, request, http.StatusOK, nil)
			return
		}

		if !methods.Contains(request.Method) {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("method not allowed")
			requests.SendResponse(respWriter, request, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("methods checked")
		next.ServeHTTP(respWriter, request)
	})
}

func RequestIDMiddleware(next http.Handler, deliver *Deliver, msg string, logger Log) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		logger.RequestID = deliver.nextRequest()
		contexted := request.WithContext(context.WithValue(context.Background(), Logg, logger))
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info(msg + " request")
		next.ServeHTTP(respWriter, contexted)
	})
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		log := request.Context().Value(Logg).(Log)
		if request.Method == http.MethodGet {
			next.ServeHTTP(respWriter, request)
			return
		}
		tok := request.Header.Get(CSRFHeader)
		if correct, err := CheckCSRFToken(request.Context().Value(RequestSID).(string),
			request.Context().Value(RequestUserID).(types.UserID), tok); !correct || err != nil {

			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("CSRF not correct")
			requests.SendResponse(respWriter, request, http.StatusForbidden, "CSRF not correct")
			return
		}
		next.ServeHTTP(respWriter, request)
	})
}
