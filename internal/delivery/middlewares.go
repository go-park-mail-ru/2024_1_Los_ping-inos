package delivery

import (
	"github.com/sirupsen/logrus"
	requests "main.go/internal/pkg"
	"net/http"
)

func IsAuthenticatedMiddleware(next http.Handler, deliver *Deliver) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		session, err := request.Cookie("session_id") // проверка авторизации
		if err != nil || session == nil || !deliver.auth.IsAuthenticated(session.Value) {
			// TODO log
			requests.SendResponse(respWriter, request, http.StatusForbidden, "forbidden")
			return
		}
		logrus.Info("check auth middleware")
		next.ServeHTTP(respWriter, request)
	})
}

func AllowedMethodMiddleware(next http.Handler, methods map[string]struct{}) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodOptions {
			logrus.Info("Preflight")
			requests.SendResponse(respWriter, request, http.StatusOK, nil)
			return
		}

		if _, ok := methods[request.Method]; !ok {
			logrus.Info(request.Method, "Method not allowed")
			requests.SendResponse(respWriter, request, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		logrus.Info("check allowed method middleware")
		next.ServeHTTP(respWriter, request)
	})
}
