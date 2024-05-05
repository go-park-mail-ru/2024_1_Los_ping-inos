package requests

import (
	"context"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	auth "main.go/internal/auth/proto"
	. "main.go/internal/logs"
	"main.go/internal/types"
	"net/http"
	"strconv"
	"time"
)

const (
	CSRFHeader = "X-Csrf-Token"
	TimingsKey = "timing"
)

func IsAuthenticatedMiddleware(next http.Handler, uc auth.AuthHandlClient) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		log := request.Context().Value(Logg).(Log)
		var session string // USER-ID

		if request.URL.Query().Has("uid") {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("authorized")
			sess, _ := strconv.Atoi(request.URL.Query().Get("uid"))
			contexted := request.WithContext(context.WithValue(request.Context(), RequestUserID, types.UserID(sess)))
			next.ServeHTTP(respWriter, contexted)
			return
		} else {
			sess, err := request.Cookie("session_id") // проверка авторизации
			if err != nil || sess == nil {
				log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("unauthorized")
				SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized")
				return
			}
			session = sess.Value
		}

		authResponse, err := uc.IsAuthenticated(request.Context(), &auth.IsAuthRequest{SessionID: session})

		if err != nil {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("unauthorized: ", err.Error())
			SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized: "+err.Error())
			return
		}

		if !authResponse.IsAuthenticated {
			log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("unauthorized")
			SendResponse(respWriter, request, http.StatusUnauthorized, "unauthorized")
			return
		}

		log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("authorized")
		contexted := request.WithContext(context.WithValue(request.Context(), RequestUserID, types.UserID(authResponse.UserID)))
		contexted = request.WithContext(context.WithValue(contexted.Context(), RequestSID, session))
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

func MetricTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		ctx = context.WithValue(ctx, TimingsKey, &ctxTimings{
			Data: make(map[string]*Timing),
		})
		// TODO log?
		defer LogContextTimings(ctx, request.URL.Path, time.Now()) // тут можно, у нас в апи нет переменных в пути
		next.ServeHTTP(respWriter, request.WithContext(ctx))
	})
}
