package service

import (
	"errors"
	"github.com/sirupsen/logrus"
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

func NewAuthHandler(dbReader PersonStorage) *AuthHandler {
	return &AuthHandler{
		sessions: make(map[string]types.UserID),
		dbReader: dbReader,
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
	person, err := api.dbReader.Get(&models.PersonGetFilter{SessionID: sessions})
	if err != nil || len(person) == 0 {
		return false
	}

	api.mutex.Lock()
	api.sessions[sessionID] = person[0].ID // нашли - запоминаем в кеш, gonka, sessii ne doljni hranitsa v sql
	api.mutex.Unlock()
	return true
}

// Login - принимает email, пароль; возвращает ID сессии и ошибку
func (api *AuthHandler) Login(email, password string) (string, error) {
	ems := make([]string, 1)
	ems[0] = email
	users, ok := api.dbReader.Get(&models.PersonGetFilter{Email: ems})
	if ok != nil {
		return "", ok
	}

	if len(users) == 0 {
		return "", errors.New("no such person")
	}

	user := users[0]

	err := checkPassword(user.Password, password)

	if err != nil {
		logrus.Info(err.Error())
		return "", errors.New("wrong password")
	}

	SID := uuid.NewString()

	api.mutex.Lock()
	api.sessions[SID] = user.ID
	user.SessionID = SID
	err = api.dbReader.Update(*user)
	if err != nil {
		logrus.Info(err.Error())
		return "", errors.New("can't write session to bd")
	}
	api.mutex.Unlock()

	return SID, nil
}

func (api *AuthHandler) Registration(Name string, Birthday string, Gender string, Email string, Password string) (string, error) {
	hashPassword, err := hashPassword(Password)
	if err != nil {
		return "", errors.New("hash func error")
	}

	err = api.dbReader.AddAccount(Name, Birthday, Gender, Email, hashPassword)
	if err != nil {
		return "", err
	}

	SID, err := api.Login(Email, Password)
	if err != nil {
		logrus.Info(err.Error())
		return "", err
	}
	return SID, nil
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

// CheckPassword - принимает hash - захэшированный пароль из базы и проверяет, соответствует ли ему password
func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
