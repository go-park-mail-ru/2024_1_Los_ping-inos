package delivery

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	requests "main.go/internal/pkg"
)

func (deliver *Deliver) GetCardsHandler(mux *http.ServeMux) {
	mux.HandleFunc("/cards",
		func(respWriter http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodOptions {
				logrus.Info("Preflight request cards")
				requests.SendResponse(respWriter, request, http.StatusOK, nil)
				return
			}

			if request.Method != http.MethodGet {
				requests.SendResponse(respWriter, request, http.StatusMethodNotAllowed, "method not allowed")
				logrus.Info("wrong method")
				return
			}

			session, err := request.Cookie("session_id") // проверка авторизации
			if err != nil || session == nil || !deliver.auth.IsAuthenticated(session.Value) {
				requests.SendResponse(respWriter, request, http.StatusForbidden, "forbidden")
				return
			}

			cards, err := deliver.serv.GetCards(session.Value)
			if err != nil {
				requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
				return
			}

			requests.SendResponse(respWriter, request, http.StatusOK, cards)
			logrus.Info("sent cards okok")
		})
}

func generateCookie(name, value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
	}
}

func setLoginCookie(sessionID, name string, expires time.Time, writer http.ResponseWriter) {
	cookie := generateCookie("session_id", sessionID, expires)
	http.SetCookie(writer, cookie)
	cookie = generateCookie("name", name, expires)
	http.SetCookie(writer, cookie)
}
