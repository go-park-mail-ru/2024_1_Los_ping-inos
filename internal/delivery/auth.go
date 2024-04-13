package delivery

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	. "main.go/internal/types"
	"net/http"
	"time"
)

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
func (deliver *Deliver) IsAuthenticatedHandler() func(w http.ResponseWriter, r *http.Request) {
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
		UID, ok := deliver.auth.IsAuthenticated(session.Value, request.Context())
		if !ok {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("not authorized")
			requests.SendResponse(respWriter, request, http.StatusUnauthorized, nil)
			return
		}

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("authorized")
		tok, err := CreateCSRFToken(session.Value, UID, oneDayExpiration().Unix())
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
func (deliver *Deliver) LoginHandler() func(respWriter http.ResponseWriter, request *http.Request) {
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

		SID, UID, err := deliver.auth.Login(request.Email, request.Password, r.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't login: ", err.Error())
			requests.SendResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		setLoginCookie(SID, oneDayExpiration(), w)
		tok, err := CreateCSRFToken(SID, UID, oneDayExpiration().Unix())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't generate csrf token: ", err.Error())
			requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("csrft", tok)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("login with SID: ", SID)
		requests.SendResponse(w, r, http.StatusOK, requests.CSRFTokenResponse{Csrft: tok})
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
func (deliver *Deliver) RegistrationHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(Logg).(*Log)
		if r.Method == http.MethodGet {
			body, err := deliver.serv.GetAllInterests(r.Context())
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
		SID, UID, err := deliver.auth.Registration(request.Name, request.Birthday, request.Gender, request.Email, request.Password, r.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't auth: ", err.Error())
			if errors.As(err, &SeveralEmailsError) {
				requests.SendResponse(w, r, http.StatusConflict, err.Error())
			} else {
				requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			}
			return
		}
		//TODO
		err = deliver.serv.UpdateProfile(SID, requests.ProfileUpdateRequest{Interests: request.Interests}, r.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("can't update interests: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		setLoginCookie(SID, oneDayExpiration(), w)
		tok, err := CreateCSRFToken(SID, UID, oneDayExpiration().Unix())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't generate csrf token: ", err.Error())
			requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("csrft", tok)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("registered and logged with SID ", SID)
		requests.SendResponse(w, r, http.StatusOK, requests.CSRFTokenResponse{Csrft: tok})
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
func (deliver *Deliver) LogoutHandler() func(respWriter http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(Logg).(*Log)

		session, err := r.Cookie("session_id")
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("no cookie 0-0 ", err.Error())
			requests.SendResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		err = deliver.auth.Logout(session.Value, r.Context())
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
