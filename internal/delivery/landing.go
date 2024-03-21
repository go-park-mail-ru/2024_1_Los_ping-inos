package delivery

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	requests "main.go/internal/pkg"
)

// GetUsername godoc
// @Summary Получить имя пользователя по его session_id (для отображения в ленте)
// @Tags Продукт
// @Router  /me [get]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Success 200		  {string}  string
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
// @Failure 500       {string} string
func (deliver *Deliver) GetUsername() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		session, err := request.Cookie("session_id")
		if err != nil || session == nil { // не знаю, нужно ли, ведь мы же сюда без куки не можем попасть :hmm:
			logrus.Info("smthg strange with cookie happened")
			requests.SendResponse(respWriter, request, http.StatusForbidden, "authenticated without cookie????")
			return
		}

		name, err := deliver.serv.GetName(session.Value)
		if err != nil {
			logrus.Info("error getting name for feed")
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, "can't get name")
		}
		requests.SendResponse(respWriter, request, http.StatusOK, name)

	}
}

// GetCardsHandler godoc
// @Summary Получить ленту
// @Tags    Продукт
// @Router  /cards [get]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Success 200		  {array}  models.Person
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
// @Failure 500       {string} string
func (deliver *Deliver) GetCardsHandler() func(http.ResponseWriter, *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		session, err := request.Cookie("session_id")
		if err != nil || session == nil { // не знаю, нужно ли, ведь мы же сюда без куки не можем попасть :hmm:
			// TODO log
			requests.SendResponse(respWriter, request, http.StatusForbidden, "authenticated without cookie????")
			return
		}

		cards, err := deliver.serv.GetCards(session.Value)
		if err != nil {
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, cards)
		logrus.Info("sent cards okok")
	}
}

func generateCookie(name, value string, expires time.Time, httpOnly bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: httpOnly,
	}
}
