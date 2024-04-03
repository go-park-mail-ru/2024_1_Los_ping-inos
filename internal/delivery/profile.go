package delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"net/http"
)

func (deliver *Deliver) ProfileHandlers() func(http.ResponseWriter, *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		switch method := request.Method; method {
		case http.MethodGet:
			deliver.ReadProfile(respWriter, request)
		case http.MethodPost:
			deliver.UpdateProfile(respWriter, request)
		case http.MethodDelete:
			deliver.DeleteProfile(respWriter, request)
		}
	}
}

// ReadProfile godoc
// @Summary Получить профиль пользователя
// @Tags    Профиль
// @Router  /profile [get]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Success 200		  {object}  models.PersonWithInterests
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
func (deliver *Deliver) ReadProfile(respWriter http.ResponseWriter, request *http.Request) {
	requestID := deliver.nextRequest()
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("read profile request")

	session, _ := request.Cookie("session_id")

	profile, err := deliver.serv.GetProfile(session.Value, requestID)
	if err != nil {
		requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("get profile err: ", err.Error())
		return
	}

	requests.SendResponse(respWriter, request, http.StatusOK, profile)
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("get profile sent response")
}

// UpdateProfile godoc
// @Summary Обновить профиль пользователя (несколько полей)
// @Description АХТУНГ АХТУНГ дата рождения передаётся в формате MM.DD.YYYY
// @Tags    Профиль
// @Router  /profile [post]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Param   userData  formData requests.ProfileUpdateRequest true "user data"
// @Success 200
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
// @Failure 409       {string} string // TODO
func (deliver *Deliver) UpdateProfile(respWriter http.ResponseWriter, request *http.Request) {
	requestID := deliver.nextRequest()
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("update profile request")

	var requestBody requests.ProfileUpdateRequest

	body, err := io.ReadAll(request.Body)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("bad body: ", err.Error())
		requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't unmarshal body: ", err.Error())
		requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	session, _ := request.Cookie("session_id")
	err = deliver.serv.UpdateProfile(session.Value, requestBody.Name, requestBody.Email, requestBody.Password, requestBody.Description,
		requestBody.Birthday, requestBody.Interests, requestID)

	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't update profile: ", err.Error())
		requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	requests.SendResponse(respWriter, request, http.StatusOK, nil)
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("update profile sent response")
}

// DeleteProfile godoc
// @Summary Удалить профиль пользователя
// @Tags    Профиль
// @Router  /profile [delete]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Success 200
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
// @Failure 409       {string} string // TODO
func (deliver *Deliver) DeleteProfile(respWriter http.ResponseWriter, request *http.Request) {
	requestID := deliver.nextRequest()
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("delete profile request")

	session, _ := request.Cookie("session_id")
	err := deliver.serv.DeleteProfile(session.Value, requestID)

	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't delete: ", err.Error())
		requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	setLoginCookie("", expiredYear, respWriter)
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("deleted profile")
	requests.SendResponse(respWriter, request, http.StatusOK, nil)
}
