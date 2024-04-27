package requests

import (
	"context"
	"encoding/json"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	. "main.go/config"
	auth "main.go/internal/auth/proto"
	. "main.go/internal/logs"
	"main.go/internal/types"
	"net/http"
)

const CSRFHeader = "csrft"

func IsAuthenticatedMiddleware(next http.Handler, _ auth.AuthHandlClient) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		log := request.Context().Value(Logg).(Log)

		session, err := request.Cookie("session_id") // проверка авторизации
		if err != nil || session == nil {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("unauthorized")
			SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized")
			return
		}

		req, err := http.NewRequest("GET", "http://127.0.0.1:8081/api/v1/isAuth", nil)
		if err != nil {
			println(err.Error())
			return
		}
		req.AddCookie(session)
		client := http.Client{}
		authResponse, _ := client.Do(req)
		body, _ := io.ReadAll(authResponse.Body)
		tmp := make(map[string]interface{}, 2)
		err = json.Unmarshal(body, &tmp)
		if _, ok := tmp["csrft"]; !ok || err != nil {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("unauthorized: ", err.Error())
			SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized: "+err.Error())
			return
		}

		log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("authorized")
		contexted := request.WithContext(context.WithValue(request.Context(), RequestUserID, types.UserID(tmp["UID"].(float64))))
		contexted = request.WithContext(context.WithValue(contexted.Context(), RequestSID, session.Value))
		next.ServeHTTP(respWriter, contexted)
	})
}

func AllowedMethodMiddleware(next http.Handler, methods *hashset.Set) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		log := request.Context().Value(Logg).(Log)

		if request.Method == http.MethodOptions {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("preflight")
			SendResponse(respWriter, request, http.StatusOK, nil)
			return
		}

		if !methods.Contains(request.Method) {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("method not allowed")
			SendResponse(respWriter, request, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("methods checked")
		next.ServeHTTP(respWriter, request)
	})
}

func RequestIDMiddleware(next http.Handler, msg string, logger Log) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		newID, _ := uuid.NewV7()
		logger.RequestID = int64(newID.ID())
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
			SendResponse(respWriter, request, http.StatusForbidden, "CSRF not correct")
			return
		}
		next.ServeHTTP(respWriter, request)
	})
}
