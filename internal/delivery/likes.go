package delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	. "main.go/config"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	"net/http"
)

// CreateLike godoc
// @Summary Создать лайк
// @Tags    Лайк
// @Router  /like [post]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Param   profile2   body   string false "profile id to like"
// @Success 200
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
func (deliver *Deliver) CreateLike() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)

		body, err := io.ReadAll(request.Body)
		if err != nil { // TODO эти два блока вынести в отдельную функцию и напихать её во все ручки
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		var requestBody requests.CreateLikeRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		err = deliver.serv.CreateLike(request.Context().Value(RequestUserID).(types.UserID), requestBody.Profile2, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't update profile: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, nil)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("create like sent response")
	}
}

// GetMatches godoc
// @Summary Получить список метчей
// @Tags    Лайк
// @Router  /matches [get]
// @Accept  json
// @Param   session_id header   string false "cookie session_id"
// @Success 200		   {array}  models.PersonWithInterests
// @Failure 400        {string} string
// @Failure 401        {string} string
// @Failure 405        {string} string
func (deliver *Deliver) GetMatches() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		userID := request.Context().Value(RequestUserID).(types.UserID)
		matches, err := deliver.serv.GetMatches(userID, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't get matches: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, matches)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("get matches sent response")
	}
}
