package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
	"strconv"

)

func (deliver *Deliver) GetCardsHandler(mux *http.ServeMux) {
	mux.HandleFunc("/",
		func(respWriter http.ResponseWriter, request *http.Request) {

			session, err := request.Cookie("session_id") // проверка авторизации
			if err != nil || session == nil || !deliver.auth.IsAuthenticated(session.Value) {
				http.Error(respWriter, "forbidden", http.StatusForbidden)
        return
			} 

			var lastID int
			if request.URL.Query().Get("last") != "" {
				lastID, err = strconv.Atoi(request.URL.Query().Get("last"))
			} else {
				lastID = 0
			}

			if err != nil {
				http.Error(respWriter, "can't return cards: no last ID", http.StatusInternalServerError)
				return
			}

			cards, err := deliver.serv.GetCards(session.Value, types.UserID(lastID))
			if err != nil {
				http.Error(respWriter, "can't return cards: smth went wrong", http.StatusInternalServerError)
				return
			}

			respWriter.Header().Set("Content-Type", "application/json")
			_, err = respWriter.Write([]byte(cards))
			if err != nil {
				http.Error(respWriter, "can't return cards: response writer error", http.StatusInternalServerError)
		})
}

func (deliver *Deliver) GetLoginHandler(mux *http.ServeMux) {
	mux.HandleFunc("/login",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				requests.SendResponse(w, r, http.StatusMethodNotAllowed, nil) // 405
				// logger
				return
			}

			var request requests.LoginRequest

			body, err := io.ReadAll(r.Body)
			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, nil)
				// logger
				return
			}

			err = json.Unmarshal(body, &request)
			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, nil)
				// logger
				return
			}

			SID, err := deliver.auth.Login(request.Email, request.Password)
			if err != nil {
				requests.SendResponse(w, r, http.StatusUnauthorized, nil)
				// logger
				return
			}

			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    SID,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true, // tolko back izmenyaet
			}

			http.SetCookie(w, cookie)

			requests.SendResponse(w, r, http.StatusOK, nil)
		})
}

func (deliver *Deliver) GetRegistrationHandler(mux *http.ServeMux) {
	mux.HandleFunc("/registration",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				requests.SendResponse(w, r, http.StatusMethodNotAllowed, nil) // 405
				// logger
				return
			}

			var request requests.RegistrationRequest

			body, err := io.ReadAll(r.Body)
			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, nil)
				// logger
				return
			}

			err = json.Unmarshal(body, &request)
			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, nil)
				// logger
				return
			}
			err = deliver.auth.Registration(request.Name, request.Birthday, request.Gender, request.Email, request.Password)

			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, nil)
			}
			requests.SendResponse(w, r, http.StatusOK, nil)
		})
}

func (deliver *Deliver) GetLogoutHandler(mux *http.ServeMux) {
	mux.HandleFunc("/logout",
		func(w http.ResponseWriter, r *http.Request) {
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

			deliver.auth.Logout(session.Value)
			if err != nil {
				requests.SendResponse(w, r, http.StatusBadRequest, nil)
				// logger
				return
			}
			requests.SendResponse(w, r, http.StatusOK, nil)
		})
}
