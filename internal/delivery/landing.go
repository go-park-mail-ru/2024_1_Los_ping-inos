package delivery

import (
	"github.com/sirupsen/logrus"
	. "main.go/config"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	"net/http"
	"time"
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
		logger := request.Context().Value(Logg).(Log)

		name, err := deliver.serv.GetName(request.Context().Value(RequestSID).(string), request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, "can't get name")
		}
		requests.SendResponse(respWriter, request, http.StatusOK, name)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent username")
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
		logger := request.Context().Value(Logg).(Log)

		cards, err := deliver.serv.GetCards(request.Context().Value(RequestUserID).(types.UserID), request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, cards)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent cards")
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
