package delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	requests "main.go/internal/pkg"
	"net/http"
	"time"
)

func setLoginCookie(sessionID, name string, expires time.Time, writer http.ResponseWriter) {
	cookie := generateCookie("session_id", sessionID, expires, true)
	http.SetCookie(writer, cookie)
	cookie = generateCookie("name", name, expires, false)
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
		session, err := request.Cookie("session_id") // проверка авторизации
		if err != nil || session == nil || !deliver.auth.IsAuthenticated(session.Value) {
			requests.SendResponse(respWriter, request, http.StatusForbidden, nil)
			return
		}
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
		var request requests.LoginRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			logrus.Info("bad request")
			return
		}

		err = json.Unmarshal(body, &request)
		if err != nil {
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			logrus.Info("can't unmarshall")
			return
		}

		SID, userName, err := deliver.auth.Login(request.Email, request.Password)
		logrus.Info("landing SID: ", SID)
		if err != nil {
			requests.SendResponse(w, r, http.StatusUnauthorized, err.Error())
			logrus.Info(err.Error())
			return
		}

		setLoginCookie(SID, userName, oneDayExpiration, w)

		requests.SendResponse(w, r, http.StatusOK, nil)
		logrus.Info("logined")
	}
}

// RegistrationHandler godoc
// @Summary Регистрация нового пользователя
// @Tags    Регистрация
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
		if r.Method == http.MethodGet {
			body, err := deliver.serv.GetAllInterests()
			if err != nil {
				requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			requests.SendResponse(w, r, http.StatusOK, body)
			return
		}

		var request requests.RegistrationRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			logrus.Info("bad request")
			return
		}

		err = json.Unmarshal(body, &request)
		if err != nil {
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			logrus.Info("can't unmarshall")
			return
		}
		SID, userName, err := deliver.auth.Registration(request.Name, request.Birthday, request.Gender, request.Email, request.Password)
		if err != nil {
			requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
			logrus.Info("can't auth")
		}

		setLoginCookie(SID, userName, oneDayExpiration, w)

		requests.SendResponse(w, r, http.StatusOK, nil)
		logrus.Info("okok")
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

		session, err := r.Cookie("session_id")
		if err != nil {
			requests.SendResponse(w, r, http.StatusUnauthorized, nil)
			logrus.Info("no cookie")
			return
		}

		err = deliver.auth.Logout(session.Value)
		if err != nil {
			requests.SendResponse(w, r, http.StatusBadRequest, nil)
			logrus.Info("can't logout")
			return
		}

		setLoginCookie("", "", expiredYear, w)

		requests.SendResponse(w, r, http.StatusOK, nil)
	}
}
