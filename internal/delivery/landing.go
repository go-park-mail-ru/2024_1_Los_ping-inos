package delivery

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	. "main.go/internal/logs"
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
		requestID := request.Context().Value(RequestID).(int64)

		session, _ := request.Cookie("session_id") // возвращает только ErrNoCookie, так что обработка не нужна

		name, err := deliver.serv.GetName(session.Value, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, "can't get name")
		}
		requests.SendResponse(respWriter, request, http.StatusOK, name)
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("sent username")
	}
}

// GetCardsHandler godoc
// @Summary Получить ленту
// @Tags    Продукт
// @Router  /cards [get]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Success 200		  {array}  models.PersonWithInterests
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
// @Failure 500       {string} string
func (deliver *Deliver) GetCardsHandler() func(http.ResponseWriter, *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		requestID := request.Context().Value(RequestID).(int64)

		session, _ := request.Cookie("session_id") // возвращает только ErrNoCookie, так что обработка не нужна

		cards, err := deliver.serv.GetCards(session.Value, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, cards)
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("sent cards")
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
