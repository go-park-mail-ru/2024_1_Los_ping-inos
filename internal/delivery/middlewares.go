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
		session, err := request.Cookie("session_id") // проверка авторизации
		if err != nil || session == nil {
			Log.WithFields(logrus.Fields{RequestID: request.Context().Value(RequestID)}).Info("unauthorized")
			requests.SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized")
			return
		}
		id, authorized := deliver.auth.IsAuthenticated(session.Value, request.Context().Value(RequestID).(int64))

		if !authorized {
			Log.WithFields(logrus.Fields{RequestID: request.Context().Value(RequestID)}).Info("unauthorized")
			requests.SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized")
			return
		}
		Log.WithFields(logrus.Fields{RequestID: request.Context().Value(RequestID)}).Info("authorized")
		contexted := request.WithContext(context.WithValue(request.Context(), RequestUserID, id))
		contexted = request.WithContext(context.WithValue(contexted.Context(), RequestSID, session.Value))
		next.ServeHTTP(respWriter, contexted)
	})
}

func AllowedMethodMiddleware(next http.Handler, methods *hashset.Set) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodOptions {
			Log.WithFields(logrus.Fields{RequestID: request.Context().Value(RequestID)}).Info("preflight")
			requests.SendResponse(respWriter, request, http.StatusOK, nil)
			return
		}

		if !methods.Contains(request.Method) {
			Log.WithFields(logrus.Fields{RequestID: request.Context().Value(RequestID)}).Info("method not allowed")
			requests.SendResponse(respWriter, request, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		Log.WithFields(logrus.Fields{RequestID: request.Context().Value(RequestID)}).Info("methods checked")
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

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			next.ServeHTTP(respWriter, request)
			return
		}
		tok := request.Header.Get(CSRFHeader)
		if correct, err := CheckCSRFToken(request.Context().Value(RequestSID).(string),
			request.Context().Value(RequestUserID).(types.UserID), tok); !correct || err != nil {

			Log.WithFields(logrus.Fields{RequestID: request.Context().Value(RequestID)}).Info("CSRF not correct")
			requests.SendResponse(respWriter, request, http.StatusForbidden, "CSRF not correct")
			return
		}
		next.ServeHTTP(respWriter, request)
	})
}
