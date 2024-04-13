package service

import (
	"context"
	"errors"
	"main.go/internal/types"
	"sync"

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

func (api *AuthHandler) IsAuthenticated(sessionID string, ctx context.Context) (types.UserID, bool) {
	person, err := api.dbReader.Get(ctx, &models.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil || len(person) == 0 {
		return -1, false
	}
	return person[0].ID, true
}

// Login - принимает email, пароль; возвращает ID сессии и ошибку
func (api *AuthHandler) Login(email, password string, ctx context.Context) (string, types.UserID, error) {
	ems := make([]string, 1)
	ems[0] = email
	users, ok := api.dbReader.Get(ctx, &models.PersonGetFilter{Email: ems})
	if ok != nil {
		return "", -1, ok
	}

	if len(users) == 0 {
		return "", -1, errors.New("no such person")
	}

	user := users[0]
	err := checkPassword(user.Password, password)

	if err != nil {
		return "", -1, err
	}

	SID := uuid.NewString()
	user.SessionID = SID
	err = api.dbReader.Update(ctx, *user)
	if err != nil {
		return "", -1, err
	}

	return SID, user.ID, nil
}

func (api *AuthHandler) Registration(name string, birthday string, gender string, email string, password string, ctx context.Context) (string, types.UserID, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", -1, err
	}

	err = api.dbReader.AddAccount(ctx, name, birthday, gender, email, hashedPassword)
	if err != nil {
		return "", -1, err
	}

	SID, UID, err := api.Login(email, password, ctx)
	if err != nil {
		return "", -1, err
	}
	return SID, UID, nil
}

func (api *AuthHandler) Logout(sessionID string, ctx context.Context) error {
	err := api.dbReader.RemoveSession(ctx, sessionID)
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
