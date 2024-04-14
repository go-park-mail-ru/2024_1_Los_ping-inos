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
	. "main.go/internal/types"
	"net/http"
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

func setLoginCookie(sessionID string, expires time.Time, writer http.ResponseWriter) {
	cookie := generateCookie("session_id", sessionID, expires, true)
	http.SetCookie(writer, cookie)
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
		logger := request.Context().Value(Logg).(*Log)
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
		UID, ok := deliver.UseCase.IsAuthenticated(session.Value, request.Context())
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
		logger := r.Context().Value(Logg).(*Log)
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
		logger := r.Context().Value(Logg).(*Log)
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
		logger := r.Context().Value(Logg).(*Log)

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

func generateCookie(name, value string, expires time.Time, httpOnly bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: httpOnly,
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
		logger := request.Context().Value(Logg).(*Log)

		name, err := deliver.UseCase.GetName(request.Context().Value(RequestSID).(string), request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, "can't get name")
		}
		requests.SendResponse(respWriter, request, http.StatusOK, name)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent username")
	}
}

func oneDayExpiration() time.Time { return time.Now().Add(24 * time.Hour) }

var expiredYear = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
