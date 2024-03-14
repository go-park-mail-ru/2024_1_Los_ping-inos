package service

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	models "main.go/db"
	"main.go/internal/types"
)

type AuthHandler struct {
	sessions map[string]types.UserID
	dbReader PersonStorage
	mutex    *sync.RWMutex
}

func NewAuthHandler(dbReader PersonStorage) *AuthHandler {
	return &AuthHandler{
		sessions: make(map[string]types.UserID),
		dbReader: dbReader,
		mutex:    &sync.RWMutex{},
	}
}

func (api *AuthHandler) IsAuthenticated(sessionID string) bool {
	// api.mutex.RLock()
	if _, authorized := api.sessions[sessionID]; authorized { // смотрим, есть ли запись в кеше
		return true
	}
	// api.mutex.RUnlock()

	// если сейчас в кеше сессии нет, лезем смотреть в бд
	sessions := make([]string, 1)
	sessions[0] = sessionID
	person, err := api.dbReader.Get(&models.PersonGetFilter{SessionID: sessions})
	if err != nil || len(person) == 0 {
		return false
	}

	// api.mutex.Lock()
	api.sessions[sessionID] = person[0].ID // нашли - запоминаем в кеш, gonka, sessii ne doljni hranitsa v sql
	// api.mutex.Unlock()
	return true
}

// Login - принимает email, пароль; возвращает ID сессии и ошибку
func (api *AuthHandler) Login(email, password string) (string, string, error) {
	ems := make([]string, 1)
	ems[0] = email
	users, ok := api.dbReader.Get(&models.PersonGetFilter{Email: ems})
	if ok != nil {
		return "", "", ok
	}

	if len(users) == 0 {
		return "", "", errors.New("no such person")
	}

	user := users[0]
	logrus.Info("LOGIN USER ", user.Email)
	err := checkPassword(user.Password, password)

	if err != nil {
		logrus.Info(err.Error())
		return "", "", errors.New("wrong password")
	}

	SID := uuid.NewString()
	logrus.Info("SID ", SID)
	api.sessions[SID] = user.ID
	user.SessionID = SID
	err = api.dbReader.Update(*user)
	logrus.Info("UPDATED")
	if err != nil {
		logrus.Info(err.Error())
		return "", "", errors.New("can't write session to bd")
	}

	return SID, user.Name, nil
}

func (api *AuthHandler) Registration(Name string, Birthday string, Gender string, Email string, Password string) (string, string, error) {
	hashPassword, err := hashPassword(Password)
	if err != nil {
		return "", "", errors.New("hash func error")
	}

	err = api.dbReader.AddAccount(Name, Birthday, Gender, Email, hashPassword)
	if err != nil {
		return "", "", err
	}

	SID, userName, err := api.Login(Email, Password)
	if err != nil {
		logrus.Info(err.Error())
		return "", "", err
	}
	return SID, userName, nil
}

func (api *AuthHandler) Logout(sessionID string) error {
	// api.mutex.RLock()
	if _, ok := api.sessions[sessionID]; !ok {
		return errors.New("no session")
	}
	// api.mutex.RUnlock()

	// api.mutex.Lock()
	delete(api.sessions, sessionID)
	// api.mutex.Unlock()

	err := api.dbReader.RemoveSession(sessionID)
	if err != nil {
		return nil
	}

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
