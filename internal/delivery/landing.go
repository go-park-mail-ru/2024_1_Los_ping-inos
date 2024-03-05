package delivery

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	"net/http"
	"strconv"
	"time"
)

func (deliver *Deliver) GetCardsHandler(mux *http.ServeMux) {
	mux.HandleFunc("/",
		func(respWriter http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodOptions {
				requests.SendResponse(respWriter, request, http.StatusOK, nil)
				return
			}

			if request.Method != http.MethodGet {
				requests.SendResponse(respWriter, request, http.StatusMethodNotAllowed, nil)
				logrus.Info("wrong method")
				return
			}

			session, err := request.Cookie("session_id") // проверка авторизации
			if err != nil || session == nil || !deliver.auth.IsAuthenticated(session.Value) {
				requests.SendResponse(respWriter, request, http.StatusForbidden, nil)
				return
			}

			var lastID int
			if request.URL.Query().Get("last") != "" {
				lastID, err = strconv.Atoi(request.URL.Query().Get("last"))
				if err != nil {
					logrus.Info("can't process ID")
					requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
				}
			} else {
				lastID = 0
			}

			cards, err := deliver.serv.GetCards(session.Value, types.UserID(lastID))
			if err != nil {
				requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
				return
			}

			_, err = respWriter.Write([]byte(cards))
			if err != nil {
				requests.SendResponse(respWriter, request, http.StatusInternalServerError,
					"can't return cards: smth went wrong")
			}
		})
}

func (deliver *Deliver) GetLoginHandler(mux *http.ServeMux) {
	mux.HandleFunc("/login",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				requests.SendResponse(w, r, http.StatusOK, nil)
				return
			}

			if r.Method != http.MethodPost {
				requests.SendResponse(w, r, http.StatusMethodNotAllowed, nil)
				// logger
				return
			}

			var request requests.LoginRequest

			body, err := io.ReadAll(r.Body)
			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
				// logger
				return
			}

			err = json.Unmarshal(body, &request)
			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
				// logger
				return
			}

			SID, err := deliver.auth.Login(request.Email, request.Password)
			if err != nil {
				requests.SendResponse(w, r, http.StatusUnauthorized, err.Error())
				logrus.Info(err.Error())
				return
			}

			cookie := generateCookie(SID)
			http.SetCookie(w, cookie)

			requests.SendResponse(w, r, http.StatusOK, nil)
		})
}

func (deliver *Deliver) GetRegistrationHandler(mux *http.ServeMux) {
	mux.HandleFunc("/registration",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				requests.SendResponse(w, r, http.StatusOK, nil)
				return
			}

			if r.Method == http.MethodGet {
				body, err := deliver.serv.GetAllInterests()
				if err != nil {
					requests.SendResponse(w, r, http.StatusInternalServerError, err.Error())
					return
				}
				requests.SendResponse(w, r, http.StatusOK, body)
				return
			}

			if r.Method != http.MethodPost {
				requests.SendResponse(w, r, http.StatusMethodNotAllowed, nil)
				logrus.Info("method not allowed")
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
			SID, err := deliver.auth.Registration(request.Name, request.Birthday, request.Gender, request.Email, request.Password)
			cookie := generateCookie(SID)
			http.SetCookie(w, cookie)

			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, err.Error())
				logrus.Info("can't auth")
			}
			requests.SendResponse(w, r, http.StatusOK, nil)
		})
}

func generateCookie(SID string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
}

func (deliver *Deliver) GetLogoutHandler(mux *http.ServeMux) {
	mux.HandleFunc("/logout",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				requests.SendResponse(w, r, http.StatusOK, nil) // 405
			}

			if r.Method != http.MethodGet { // delete zapros?
				requests.SendResponse(w, r, http.StatusMethodNotAllowed, nil) // 405
				// logger
				return
			}
			session, err := r.Cookie("session_id")
			if err != nil {
				requests.SendResponse(w, r, http.StatusUnauthorized, nil)
				// logger
				return
			}

			err = deliver.auth.Logout(session.Value)
			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, nil)
				// logger
				return
			}
			requests.SendResponse(w, r, http.StatusOK, nil)
		})
}
