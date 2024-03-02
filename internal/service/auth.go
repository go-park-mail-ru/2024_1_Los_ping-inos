package service

import (
	"main.go/db"
	"main.go/internal/types"
	"math/rand"
	"net/http"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type Auth interface {
	IsAuthenticated(w http.ResponseWriter, r *http.Request) bool
	Login(w http.ResponseWriter, r *http.Request) bool
	Logout(w http.ResponseWriter, r *http.Request) bool
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type AuthHandler struct {
	sessions map[string]types.UserID
	dbReader Storage
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		sessions: make(map[string]types.UserID),
	}
}

func (api *AuthHandler) IsAuthenticated(_ http.ResponseWriter, r *http.Request) bool {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		_, authorized = api.sessions[session.Value]
	}

	if authorized {
		return true
	} else { // если сейчас в кеше сессии нет, лезем смотреть в бд
		sessions := make([]string, 1)
		sessions[0] = session.Value
		person, err := api.dbReader.Get(&models.PersonFilter{SessionID: sessions})
		if err != nil || person == nil {
			return false
		}
		api.sessions[session.Value] = person[0].ID // нашли - запоминаем
		return true
	}
}

func (api *AuthHandler) Login(w http.ResponseWriter, r *http.Request) bool {
	ems := make([]string, 1)
	ems[0] = r.FormValue("email")
	users, ok := api.dbReader.Get(&models.PersonFilter{Email: ems})
	if ok != nil {
		http.Error(w, `no user`, http.StatusNotFound)
		return false
	}
	user := users[0]
	if user.Password != r.FormValue("password") {
		http.Error(w, `bad pass`, http.StatusBadRequest)
		return false
	}

	SID := RandStringRunes(32)

	api.sessions[SID] = user.ID

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	return true
}

func (api *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) bool {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, `no sess`, http.StatusUnauthorized)
		return false
	}

	if _, ok := api.sessions[session.Value]; !ok {
		http.Error(w, `no sess`, http.StatusUnauthorized)
		return false
	}

	delete(api.sessions, session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	return true
}
