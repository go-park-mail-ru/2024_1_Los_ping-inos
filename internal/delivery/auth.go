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
		requestID := deliver.nextRequest()
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("auth check")
		session, err := request.Cookie("session_id") // проверка авторизации

		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("not authorized: ", err.Error())
		}
		if err != nil || session == nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Info("not authorized")
			requests.SendResponse(respWriter, request, http.StatusForbidden, nil)
			return
		}
		_, ok := deliver.auth.IsAuthenticated(session.Value, requestID)
		if !ok {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Info("not authorized")
			requests.SendResponse(respWriter, request, http.StatusForbidden, nil)
			return
		}

		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("authorized")
		requests.SendResponse(respWriter, request, http.StatusOK, nil)
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
		requestID := deliver.nextRequest()
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("login")
		var request requests.LoginRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("bad body: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &request)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		SID, UID, err := deliver.auth.Login(request.Email, request.Password, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't login: ", err.Error())
			requests.SendResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		setLoginCookie(SID, oneDayExpiration(), w)
		tok, err := CreateCSRFToken(SID, UID, oneDayExpiration().Unix())
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't generate csrf: ", err.Error())
			requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("csrft", tok)
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("login with SID: ", SID)
		tmp, _ := json.Marshal(requests.CSRFTokenResponse{tok})
		requests.SendResponse(w, r, http.StatusOK, string(tmp))
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
		requestID := r.Context().Value(RequestID).(int64)
		if r.Method == http.MethodGet {
			body, err := deliver.serv.GetAllInterests(requestID)
			if err != nil {
				Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't get interests: ", err.Error())
				requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			Log.WithFields(logrus.Fields{RequestID: requestID}).Info("sent interests")
			requests.SendResponse(w, r, http.StatusOK, body)
			return
		}

		var request requests.RegistrationRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("bad body: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &request)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}
		SID, UID, err := deliver.auth.Registration(request.Name, request.Birthday, request.Gender, request.Email, request.Password, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't auth: ", err.Error())
			if errors.As(err, &SeveralEmailsError) {
				requests.SendResponse(w, r, http.StatusConflict, err.Error())
			} else {
				requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			}
			return
		}
		//TODO
		err = deliver.serv.UpdateProfile(SID, requests.ProfileUpdateRequest{Interests: request.Interests}, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Info("can't update interests: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		setLoginCookie(SID, oneDayExpiration(), w)
		tok, err := CreateCSRFToken(SID, UID, oneDayExpiration().Unix())
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't generate csrf: ", err.Error())
			requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("csrft", tok)
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("registered and logged with SID ", SID)
		tmp, _ := json.Marshal(requests.CSRFTokenResponse{tok})
		requests.SendResponse(w, r, http.StatusOK, string(tmp))
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
		requestID := r.Context().Value(RequestID).(int64)

		session, err := r.Cookie("session_id")
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("no cookie 0-0 ", err.Error())
			requests.SendResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		err = deliver.auth.Logout(session.Value, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't logout: ", err.Error())
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		setLoginCookie("", expiredYear, w)
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("logout end")
		requests.SendResponse(w, r, http.StatusOK, nil)
	}
}
