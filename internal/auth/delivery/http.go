package delivery

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	. "main.go/config"
	"main.go/internal/auth"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	. "main.go/internal/types"
	"net/http"
	"strconv"
	"time"
)

type AuthHandler struct {
	UseCase auth.IUseCase
}

func NewAuthHandler(uc auth.IUseCase) *AuthHandler {
	return &AuthHandler{
		UseCase: uc,
	}
}

func (deliver *AuthHandler) ProfileHandlers() func(http.ResponseWriter, *http.Request) {
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
// @Success 200		  {object}  auth.PersonWithInterests
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 405       {string} string
func (deliver *AuthHandler) ReadProfile(respWriter http.ResponseWriter, request *http.Request) {
	var (
		err  error
		id   int
		prof []auth.Profile
	)

	logger := request.Context().Value(Logg).(Log)

	if request.URL.Query().Has("id") { // просмотр профиля по id (чужой профиль из ленты)
		id, err = strconv.Atoi(request.URL.Query().Get("id"))
		if err != nil {
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("get profile err: ", err.Error())
		}
		prof, err = deliver.UseCase.GetProfile(auth.ProfileGetParams{ID: []types.UserID{types.UserID(id)}}, request.Context())
	} else { // свой профиль
		prof, err = deliver.UseCase.GetProfile(auth.ProfileGetParams{ID: []types.UserID{request.Context().Value(RequestUserID).(types.UserID)}}, request.Context())
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
func (deliver *AuthHandler) UpdateProfile(respWriter http.ResponseWriter, request *http.Request) {
	logger := request.Context().Value(Logg).(Log)

	var requestBody auth.ProfileUpdateRequest

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

	err = deliver.UseCase.UpdateProfile(session.Value, requestBody, request.Context())
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
func (deliver *AuthHandler) DeleteProfile(respWriter http.ResponseWriter, request *http.Request) {
	logger := request.Context().Value(Logg).(Log)

	err := deliver.UseCase.DeleteProfile(request.Context().Value("SID").(string), request.Context())

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
// @Success 200		   {array}  auth.PersonWithInterests
// @Failure 400        {string} string
// @Failure 401        {string} string
// @Failure 405        {string} string
func (deliver *AuthHandler) GetMatches() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		userID := request.Context().Value(RequestUserID).(types.UserID)
		matches, err := deliver.UseCase.GetMatches(userID, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't get matches: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, matches)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("get matches sent response")
	}
}

// IsAuthenticatedHandler godoc
// @Summary Проверка авторизации пользователя
// @Description Проверка по session_id из куки (если она есть)
// @Tags    Авторизация
// @Router  /isAuth [get]
// @Param  session_id header string false "cookie session_id"
// @Success 200
// @Failure 403
func (deliver *AuthHandler) IsAuthenticatedHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("auth check")
		session, err := request.Cookie("session_id") // проверка авторизации

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("not authorized: ", err.Error())
		}
		if err != nil || session == nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("not authorized")
			requests.SendResponse(respWriter, request, http.StatusUnauthorized, nil)
			return
		}
		UID, ok, err := deliver.UseCase.IsAuthenticated(session.Value, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info(err.Error())
			requests.SendResponse(respWriter, request, http.StatusUnauthorized, nil)
			return
		}
		if !ok {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("not authorized")
			requests.SendResponse(respWriter, request, http.StatusUnauthorized, nil)
			return
		}

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("authorized")
		tok, err := requests.CreateCSRFToken(session.Value, UID, oneDayExpiration().Unix())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't generate csrf token: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}
		requests.SendResponse(respWriter, request, http.StatusOK, requests.CSRFTokenResponse{Csrft: tok})
	}
}

// LoginHandler godoc
// @Summary Залогинить пользователя
// @Tags    Авторизация
// @Router  /login [post]
// @Accept  json
// @Param   userData  formData requests.LoginRequest true "user data"
// @Success 200
// @Failure 405       {string} string
// @Failure 400       {string} string
// @Failure 401       {string} string
func (deliver *AuthHandler) LoginHandler() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(Logg).(Log)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("login")
		var request requests.LoginRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("bad body: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &request)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		prof, SID, err := deliver.UseCase.Login(request.Email, request.Password, r.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't login: ", err.Error())
			requests.SendResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		setLoginCookie(SID, oneDayExpiration(), w)
		tok, err := requests.CreateCSRFToken(SID, prof.ID, oneDayExpiration().Unix())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't generate csrf token: ", err.Error())
			requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("csrft", tok)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("login with SID: ", SID)
		prof.CSRFT = tok
		requests.SendResponse(w, r, http.StatusOK, prof)
	}
}

// RegistrationHandler godoc
// @Summary Регистрация нового пользователя
// @Description АХТУНГ АХТУНГ дата рождения передаётся в формате MM.DD.YYYY
// @Tags    Профиль
// @Router  /registration [post]
// @Accept  json
// @Param   userData  formData requests.RegistrationRequest true "user data"
// @Success 200
// @Failure 405       {string} string
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 500       {string} string
func (deliver *AuthHandler) RegistrationHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(Logg).(Log)
		if r.Method == http.MethodGet {
			body, err := deliver.UseCase.GetAllInterests(r.Context())
			if err != nil {
				logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't get interests: ", err.Error())
				requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent interests")
			requests.SendResponse(w, r, http.StatusOK, body)
			return
		}

		var request requests.RegistrationRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("bad body: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &request)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}
		prof, SID, err := deliver.UseCase.Registration(auth.RegitstrationBody{request.Name, request.Birthday,
			request.Gender, request.Email, request.Password, request.Interests}, r.Context())

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't auth: ", err.Error())
			if errors.As(err, &SeveralEmailsError) {
				requests.SendResponse(w, r, http.StatusConflict, err.Error())
			} else {
				requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			}
			return
		}
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("can't update interests: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		setLoginCookie(SID, oneDayExpiration(), w)
		tok, err := requests.CreateCSRFToken(SID, prof.ID, oneDayExpiration().Unix())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't generate csrf token: ", err.Error())
			requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("csrft", tok)
		prof.CSRFT = tok
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("registered and logged with SID ", SID)
		requests.SendResponse(w, r, http.StatusOK, prof)
	}
}

// LogoutHandler godoc
// @Summary Разлогин
// @Tags    Авторизация
// @Router  /logout [get]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Success 200
// @Failure 405       {string} string
// @Failure 400       {string} string
// @Failure 401       {string} string
// @Failure 500       {string} string
func (deliver *AuthHandler) LogoutHandler() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(Logg).(Log)

		session, err := r.Cookie("session_id")
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("no cookie 0-0 ", err.Error())
			requests.SendResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		err = deliver.UseCase.Logout(session.Value, r.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't logout: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		setLoginCookie("", expiredYear, w)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("logout end")
		requests.SendResponse(w, r, http.StatusOK, nil)
	}
}

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
func (deliver *AuthHandler) GetUsername() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)

		name, err := deliver.UseCase.GetName(request.Context().Value(RequestUserID).(types.UserID), request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, "can't get name")
		}
		requests.SendResponse(respWriter, request, http.StatusOK, name)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent username")
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

func oneDayExpiration() time.Time { return time.Now().Add(24 * time.Hour) }

var expiredYear = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
