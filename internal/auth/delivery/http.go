package delivery

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mailru/easyjson"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
	. "main.go/config"
	"main.go/internal/auth"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	. "main.go/internal/types"
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

	//logger := request.Context().Value(Logg).(Log)

	if request.URL.Query().Has("id") { // просмотр профиля по id (чужой профиль из ленты)
		id, err = strconv.Atoi(request.URL.Query().Get("id"))
		if err != nil {
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("get profile err: ", err.Error())
		}
		prof, err = deliver.UseCase.GetProfile(auth.ProfileGetParams{ID: []types.UserID{types.UserID(id)}, NeedEmail: false}, request.Context())
	} else { // свой профиль
		prof, err = deliver.UseCase.GetProfile(auth.ProfileGetParams{ID: []types.UserID{types.UserID(id)}, NeedEmail: true}, request.Context())
	}

	if err != nil {
		requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
		//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("get profile err: ", err.Error())
		return
	}

	if len(prof) == 0 {
		requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, "no such profile")
		//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("no such profile")
		return
	}

	requests.SendResponse(respWriter, request, http.StatusOK, prof[0])
	//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("get profile sent response")
}

// UpdateProfile godoc
// @Summary Обновить профиль пользователя (несколько полей)
// @Description АХТУНГ АХТУНГ дата рождения передаётся в формате MM.DD.YYYY
// @Tags    Профиль
// @Router  /profile [post]
// @Accept  json
// @Param   session_id header string false "cookie session_id"
// @Param   userData  formData auth.ProfileUpdateRequest true "user data"
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
		requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	if err = easyjson.Unmarshal(body, &requestBody); err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
		requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	err = deliver.UseCase.UpdateProfile(request.Context().Value(1).(types.UserID), requestBody, request.Context())
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't update profile: ", err.Error())
		if errors.As(err, &types.MyErr{Err: types.DifferentPasswordsError}) {
			requests.SendSimpleResponse(respWriter, request, http.StatusConflict, err.Error())
		} else {
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
		}
		return
	}

	requests.SendSimpleResponse(respWriter, request, http.StatusOK, "")
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

	err := deliver.UseCase.DeleteProfile(request.Context().Value(RequestUserID).(types.UserID), request.Context())

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't delete: ", err.Error())
		requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
		return
	}

	setLoginCookie("", expiredYear, respWriter)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("deleted profile")
	requests.SendSimpleResponse(respWriter, request, http.StatusOK, "")
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
		var requestBody auth.GetMatchesRequest
		body, err := io.ReadAll(request.Body)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("bad body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		//err = json.Unmarshal(body, &requestBody) TODO
		if err = easyjson.Unmarshal(body, &requestBody); err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		requestBody.Name = strings.ToLower(requestBody.Name)
		matches, err := deliver.UseCase.GetMatches(userID, requestBody.Name, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't get matches: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, auth.Matches{Match: matches})
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
			requests.SendSimpleResponse(respWriter, request, http.StatusUnauthorized, "")
			return
		}
		UID, ok, err := deliver.UseCase.IsAuthenticated(session.Value, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusUnauthorized, "")
			return
		}
		if !ok {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("not authorized")
			requests.SendSimpleResponse(respWriter, request, http.StatusUnauthorized, "")
			return
		}

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("authorized")
		tok, err := requests.CreateCSRFToken(session.Value, UID, oneDayExpiration().Unix())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't generate csrf token: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}
		requests.SendResponse(respWriter, request, http.StatusOK, auth.CSRFTokenResponse{Csrft: tok})
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
		var request auth.LoginRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("bad body: ", err.Error())
			requests.SendSimpleResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		//err = json.Unmarshal(body, &request) TODO
		if err = easyjson.Unmarshal(body, &request); err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendSimpleResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		prof, SID, err := deliver.UseCase.Login(request.Email, request.Password, r.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't login: ", err.Error())
			requests.SendSimpleResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		setLoginCookie(SID, oneDayExpiration(), w)
		tok, err := requests.CreateCSRFToken(SID, prof.ID, oneDayExpiration().Unix())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't generate csrf token: ", err.Error())
			requests.SendSimpleResponse(w, r, http.StatusInternalServerError, err.Error())
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
				requests.SendSimpleResponse(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent interests")
			requests.SendResponse(w, r, http.StatusOK, auth.Interests{Interes: body})
			return
		}

		var request auth.RegistrationRequest

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("bad body: ", err.Error())
			requests.SendSimpleResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		//err = json.Unmarshal(body, &request) TODO
		if err = easyjson.Unmarshal(body, &request); err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendSimpleResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}
		prof, SID, err := deliver.UseCase.Registration(auth.RegitstrationBody{Name: request.Name, Birthday: request.Birthday,
			Gender: request.Gender, Email: request.Email, Password: request.Password, Interests: request.Interests}, r.Context())

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't auth: ", err.Error())
			if errors.As(err, &types.MyErr{Err: SeveralEmailsError}) {
				requests.SendSimpleResponse(w, r, http.StatusConflict, err.Error())
			} else {
				requests.SendSimpleResponse(w, r, http.StatusBadRequest, err.Error())
			}
			return
		}
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("can't update interests: ", err.Error())
			requests.SendSimpleResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		setLoginCookie(SID, oneDayExpiration(), w)
		tok, err := requests.CreateCSRFToken(SID, prof.ID, oneDayExpiration().Unix())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't generate csrf token: ", err.Error())
			requests.SendSimpleResponse(w, r, http.StatusInternalServerError, err.Error())
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
			requests.SendSimpleResponse(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		err = deliver.UseCase.Logout(session.Value, r.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't logout: ", err.Error())
			requests.SendSimpleResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		setLoginCookie("", expiredYear, w)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("logout end")
		requests.SendSimpleResponse(w, r, http.StatusOK, "")
	}
}

func MetricTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next.ServeHTTP(respWriter, request)

		end := time.Since(start)
		path := request.URL.Path
		if path != "/metrics" {
			auth.TotalHits.WithLabelValues().Inc()
			auth.HitDuration.WithLabelValues(request.Method, path).Set(float64(end.Milliseconds()))
		}
	})
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
		Domain:   ".jimder.ru",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
}

func oneDayExpiration() time.Time { return time.Now().Add(24 * time.Hour) }

var expiredYear = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

func (deliver *AuthHandler) PaymentUrl() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		UID := request.Context().Value(RequestUserID).(types.UserID)

		urll := deliver.UseCase.GenPaymentUrl(UID)

		requests.SendSimpleResponse(respWriter, request, http.StatusOK, urll)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent paymentUrl")
	}
}

func (deliver *AuthHandler) ActivateSub() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		UID := request.Context().Value(RequestUserID).(types.UserID)

		datetime, activated := deliver.checkSubStatus(UID)
		if activated != nil && errors.As(activated, &types.MyErr{Err: auth.NoPaymentErr}) {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("payment not provided")
			requests.SendSimpleResponse(respWriter, request, http.StatusConflict, "payment not provided")
			return
		}
		if activated != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't check payment: ", activated.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusConflict, "can't check payment: "+activated.Error())
			return
		}

		err := deliver.UseCase.ActivateSub(request.Context(), UID, datetime)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't activate sub: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		exp := datetime.Add(31 * 24 * time.Hour).Unix()
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("sent response activating sub")
		requests.SendSimpleResponse(respWriter, request, http.StatusOK, strconv.Itoa(int(exp)))
		return
	}
}

func (deliver *AuthHandler) checkSubStatus(UID types.UserID) (time.Time, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "yoomoney.ru",
		Path:   "/api/operation-history",
	}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return time.Now(), err
	}
	key := viper.Get("yoomoney.key").(string)
	req.Header.Add("Authorization", "Bearer "+key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return time.Now(), err
	}
	defer resp.Body.Close()

	var operations auth.Operations

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return time.Now(), err
	}
	err = easyjson.Unmarshal(body, &operations)
	uid := strconv.Itoa(int(UID))
	for _, i := range operations.Operations {
		if i.Label == uid {
			if time.Since(i.Datetime) <= 31*24*time.Hour {
				return i.Datetime, nil
			}
		}
	}
	return time.Now(), auth.NoPaymentErr
}

func (deliver *AuthHandler) GetSubHistory() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		UID := request.Context().Value(RequestUserID).(types.UserID)

		res, err := deliver.UseCase.GetSubHistory(request.Context(), UID)

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't read sub history: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("sent history")
		requests.SendResponse(respWriter, request, http.StatusOK, res)
	}
}
