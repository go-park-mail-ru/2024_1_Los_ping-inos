package service

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	models "main.go/db"
	"main.go/internal/types"
)

type AuthHandler struct {
	sessions map[string]types.UserID
	dbReader PersonStorage
	mutex    sync.RWMutex
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		sessions: make(map[string]types.UserID),
	}
}

func (api *AuthHandler) IsAuthenticated(sessionID string) bool {
	api.mutex.RLock()
	if _, authorized := api.sessions[sessionID]; authorized { // смотрим, есть ли запись в кеше
		return true
	}
	api.mutex.RUnlock()

	// если сейчас в кеше сессии нет, лезем смотреть в бд
	sessions := make([]string, 1)
	sessions[0] = sessionID
	person, err := api.dbReader.Get(&models.PersonFilter{SessionID: sessions})
	if err != nil || person == nil {
		return false
	}

	api.mutex.Lock()
	api.sessions[sessionID] = person[0].ID // нашли - запоминаем в кеш, gonka, sessii ne doljni hranitsa v sql
	api.mutex.Unlock()
	return true
}

// Login - принимает логин, пароль; возвращает ID сессии и ошибку
// Email and password, not login
func (api *AuthHandler) Login(email, password string) (string, error) {
	ems := make([]string, 1)
	ems[0] = email
	users, ok := api.dbReader.Get(&models.PersonFilter{Email: ems})
	if ok != nil || users == nil {
		return "", errors.New("no such person")
	}

	user := users[0]

	pass, err := hashPassword(password)
	if err != nil || user.Password != pass { // dve raznie proverki doljni bit
		return "", errors.New("wrong password")
	}

	SID := uuid.NewString()

	api.mutex.Lock()
	api.sessions[SID] = user.ID // gonka
	api.mutex.Unlock()

	// TODO сделать update в personStorage и через него в бд записать

	return SID, nil
}

func (api *AuthHandler) Registration(Name string, Birthday string, Gender string, Email string, Password string) error {
	hashPassword, err := hashPassword(Password)
	if err != nil {
		return errors.New("hash func error")
	}

	err = api.dbReader.AddAccount(Name, Birthday, Gender, Email, hashPassword)
	if err != nil {
		return err
	}

	return nil
}

func (api *AuthHandler) Logout(sessionID string) error {
	api.mutex.RLock()
	if _, ok := api.sessions[sessionID]; !ok {
		return errors.New("no session")
	}
	api.mutex.RUnlock()

	api.mutex.Lock()
	delete(api.sessions, sessionID)
	api.mutex.Unlock()

	// TODO сделать update в personStorage и через него в бд записать

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // TODO подумать насчет константы
	return string(bytes), err
}
