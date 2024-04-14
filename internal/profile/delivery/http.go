package delivery

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"main.go/config"
	. "main.go/config"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/profile"
	"main.go/internal/types"
	"net/http"
	"strconv"
	"time"
)

type ProfileHandler struct {
	useCase profile.UseCase
}

func NewProfileDeliver(uc profile.UseCase) *ProfileHandler {
	return &ProfileHandler{
		useCase: uc,
	}
}

func (deliver *ProfileHandler) ProfileHandlers() func(http.ResponseWriter, *http.Request) {
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
// @Param   id         query  string false "profile id to return (optional)"
// @Success 200		  {object}  profile.PersonWithInterests
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
func (deliver *ProfileHandler) ReadProfile(respWriter http.ResponseWriter, request *http.Request) {
	var (
		err  error
		id   int
		prof []profile.Card
	)

	logger := request.Context().Value(Logg).(*Log)

	if request.URL.Query().Has("id") { // просмотр профиля по id (чужой профиль из ленты)
		id, err = strconv.Atoi(request.URL.Query().Get("id"))
		if err != nil {
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("get profile err: ", err.Error())
		}
		prof, err = deliver.useCase.GetProfile(profile.ProfileGetParams{ID: []types.UserID{types.UserID(id)}}, request.Context())
	} else { // свой профиль
		session, _ := request.Cookie("session_id")
		prof, err = deliver.useCase.GetProfile(profile.ProfileGetParams{SessionID: []string{session.Value}}, request.Context())
	}

	if err != nil {
		requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("get profile err: ", err.Error())
		return
	}

	requests.SendResponse(respWriter, request, http.StatusOK, prof[0])
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("get profile sent response")
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
func (deliver *ProfileHandler) UpdateProfile(respWriter http.ResponseWriter, request *http.Request) {
	logger := request.Context().Value(Logg).(*Log)

	var requestBody requests.ProfileUpdateRequest

	body, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
		requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
		requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	session, _ := request.Cookie("session_id")

	err = deliver.useCase.UpdateProfile(session.Value, requestBody, request.Context())
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't update profile: ", err.Error())
		if errors.As(err, &types.DifferentPasswordsError) {
			requests.SendResponse(respWriter, request, http.StatusConflict, err.Error())
		} else {
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		}
		return
	}

	requests.SendResponse(respWriter, request, http.StatusOK, nil)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("update profile sent response")
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
func (deliver *ProfileHandler) DeleteProfile(respWriter http.ResponseWriter, request *http.Request) {
	logger := request.Context().Value(Logg).(*Log)

	err := deliver.useCase.DeleteProfile(request.Context().Value(config.RequestSID).(string), request.Context())

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't delete: ", err.Error())
		requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	setLoginCookie("", expiredYear, respWriter)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("deleted profile")
	requests.SendResponse(respWriter, request, http.StatusOK, nil)
}

// GetMatches godoc
// @Summary Получить список метчей
// @Tags    Лайк
// @Router  /matches [get]
// @Accept  json
// @Param   session_id header   string false "cookie session_id"
// @Success 200		   {array}  profile.PersonWithInterests
// @Failure 400        {string} string
// @Failure 401        {string} string
// @Failure 405        {string} string
func (deliver *ProfileHandler) GetMatches() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(*Log)
		userID := request.Context().Value(RequestUserID).(types.UserID)
		matches, err := deliver.useCase.GetMatches(userID, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't get matches: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, matches)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("get matches sent response")
	}
}

func setLoginCookie(sessionID string, expires time.Time, writer http.ResponseWriter) {
	cookie := generateCookie("session_id", sessionID, expires, true)
	http.SetCookie(writer, cookie)
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

var expiredYear = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
