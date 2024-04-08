package service

import (
	"errors"
	"main.go/internal/types"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	models "main.go/db"
	. "main.go/internal/logs"
)

type AuthHandler struct {
	sessions *sync.Map
	dbReader PersonStorage
}

func NewAuthHandler(dbReader PersonStorage) *AuthHandler {
	return &AuthHandler{
		sessions: &sync.Map{},
		dbReader: dbReader,
	}
}

func (api *AuthHandler) IsAuthenticated(sessionID string, requestID int64) (types.UserID, bool) {
	if id, authorized := api.sessions.Load(sessionID); authorized { // смотрим, есть ли запись в кеше
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("loaded session ", sessionID)
		return id.(types.UserID), true
	}

	// если сейчас в кеше сессии нет, лезем смотреть в бд
	sessions := make([]string, 1)
	sessions[0] = sessionID
	person, err := api.dbReader.Get(requestID, &models.PersonGetFilter{SessionID: sessions})
	if err != nil || len(person) == 0 {
		return -1, false
	}

	api.sessions.Store(sessionID, person[0].ID) // нашли - запоминаем в кеш
	return person[0].ID, true
}

// Login - принимает email, пароль; возвращает ID сессии и ошибку
func (api *AuthHandler) Login(email, password string, requestID int64) (string, error) {
	ems := make([]string, 1)
	ems[0] = email
	users, ok := api.dbReader.Get(requestID, &models.PersonGetFilter{Email: ems})
	if ok != nil {
		return "", ok
	}

	if len(users) == 0 {
		return "", errors.New("no such person")
	}

	user := users[0]
	err := checkPassword(user.Password, password)

	if err != nil {
		return "", err
	}

	SID := uuid.NewString()
	api.sessions.Store(SID, user.ID)
	user.SessionID = SID
	err = api.dbReader.Update(requestID, *user)
	if err != nil {
		return "", err
	}

	return SID, nil
}

func (api *AuthHandler) Registration(name string, birthday string, gender string, email string, password string, requestID int64) (string, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", err
	}

	err = api.dbReader.AddAccount(requestID, name, birthday, gender, email, hashedPassword)
	if err != nil {
		return "", err
	}

	SID, err := api.Login(email, password, requestID)
	if err != nil {
		return "", err
	}
	return SID, nil
}

func (api *AuthHandler) Logout(sessionID string, requestID int64) error {
	if _, ok := api.sessions.Load(sessionID); !ok {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("no session to logout")
		return errors.New("no session")
	}

	api.sessions.Delete(sessionID)

	err := api.dbReader.RemoveSession(requestID, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword - принимает hash - захэшированный пароль из базы и проверяет, соответствует ли ему password
func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
