package delivery

import (
	requests "main.go/internal/pkg"
	"net/http"
)

func IsAuthenticatedMiddleware(next http.Handler, deliver *Deliver) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		session, err := request.Cookie("session_id")                                         // проверка авторизации
		if err != nil || session == nil || !deliver.auth.IsAuthenticated(session.Value, 0) { // я хз, как сделать один id у мидлвары и следующего обработчика
			requests.SendResponse(respWriter, request, http.StatusForbidden, "forbidden")
			return
		}
		next.ServeHTTP(respWriter, request)
	})
}

func AllowedMethodMiddleware(next http.Handler, methods map[string]struct{}) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodOptions {
			requests.SendResponse(respWriter, request, http.StatusOK, nil)
			return
		}

		if _, ok := methods[request.Method]; !ok {
			requests.SendResponse(respWriter, request, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		next.ServeHTTP(respWriter, request)
	})
}
