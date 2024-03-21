package service

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	models "main.go/db"
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

func (api *AuthHandler) IsAuthenticated(sessionID string) bool {
	if _, authorized := api.sessions.Load(sessionID); authorized { // смотрим, есть ли запись в кеше
		logrus.Info("loaded ", sessionID)
		return true
	}

	// если сейчас в кеше сессии нет, лезем смотреть в бд
	sessions := make([]string, 1)
	sessions[0] = sessionID
	person, err := api.dbReader.Get(&models.PersonGetFilter{SessionID: sessions})
	if err != nil || len(person) == 0 {
		logrus.Info("no such person")
		return false
	}

	logrus.Info("person ", person[0].Name)
	api.sessions.Store(sessionID, person[0].ID) // нашли - запоминаем в кеш, gonka, sessii ne doljni hranitsa v sql

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
	api.sessions.Store(SID, user.ID)
	user.SessionID = SID
	err = api.dbReader.Update(*user)
	logrus.Info("UPDATED")
	if err != nil {
		logrus.Info(err.Error())
		return "", "", errors.New("can't write session to bd")
	}

	return SID, user.Name, nil
}

func (api *AuthHandler) Registration(name string, birthday string, gender string, email string, password string) (string, string, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", "", errors.New("hash func error")
	}

	err = api.dbReader.AddAccount(name, birthday, gender, email, hashedPassword)
	if err != nil {
		return "", "", err
	}

	SID, userName, err := api.Login(email, password)
	if err != nil {
		logrus.Info(err.Error())
		return "", "", err
	}
	return SID, userName, nil
}

func (api *AuthHandler) Logout(sessionID string) error {
	if _, ok := api.sessions.Load(sessionID); !ok {
		return errors.New("no session")
	}

	api.sessions.Delete(sessionID)

	err := api.dbReader.RemoveSession(sessionID)
	if err != nil {
		return nil
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
