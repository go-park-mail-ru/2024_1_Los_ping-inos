package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"main.go/db"
	"main.go/internal/types"
)

type AuthHandler struct {
	sessions map[string]types.UserID
	dbReader PersonStorage
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		sessions: make(map[string]types.UserID),
	}
}

func (api *AuthHandler) IsAuthenticated(sessionID string) bool {
	if _, authorized := api.sessions[sessionID]; authorized { // смотрим, есть ли запись в кеше
		return true
	}

	// если сейчас в кеше сессии нет, лезем смотреть в бд
	sessions := make([]string, 1)
	sessions[0] = sessionID
	person, err := api.dbReader.Get(&models.PersonFilter{SessionID: sessions})
	if err != nil || person == nil {
		return false
	}

	api.sessions[sessionID] = person[0].ID // нашли - запоминаем в кеш
	return true
}

// Login - принимает логин, пароль; возвращает ID сессии и ошибку
func (api *AuthHandler) Login(email, password string) (string, error) {
	ems := make([]string, 1)
	ems[0] = email
	users, ok := api.dbReader.Get(&models.PersonFilter{Email: ems})
	if ok != nil || users == nil {
		return "", errors.New("no such person")
	}

	user := users[0]

	pass, err := hashPassword(password)
	if err != nil || user.Password != pass {
		return "", errors.New("wrong password")
	}

	SID := uuid.NewString()

	api.sessions[SID] = user.ID

	// TODO сделать update в personStorage и через него в бд записать

	return SID, nil
}

func (api *AuthHandler) Logout(sessionID string) error {
	if _, ok := api.sessions[sessionID]; !ok {
		return errors.New("no session")
	}

	delete(api.sessions, sessionID)

	// TODO сделать update в personStorage и через него в бд записать

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // TODO подумать насчет константы
	return string(bytes), err
}
