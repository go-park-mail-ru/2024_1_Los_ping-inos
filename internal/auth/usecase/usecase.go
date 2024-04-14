package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"main.go/internal/auth"
	"main.go/internal/types"
)

type UseCase struct {
	dbReader        auth.PostgresRepo
	interestStorage auth.InterestStorage
}

func NewAuthUseCase(dbReader auth.PostgresRepo, istore auth.InterestStorage) *UseCase {
	return &UseCase{
		dbReader:        dbReader,
		interestStorage: istore,
	}
}

func (api *UseCase) IsAuthenticated(sessionID string, ctx context.Context) (types.UserID, bool) {
	person, err := api.dbReader.Get(ctx, &auth.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil || len(person) == 0 {
		return -1, false
	}
	return person[0].ID, true
}

// Login - принимает email, пароль; возвращает ID сессии и ошибку
func (api *UseCase) Login(email, password string, ctx context.Context) (string, types.UserID, error) {
	ems := make([]string, 1)
	ems[0] = email
	users, ok := api.dbReader.Get(ctx, &auth.PersonGetFilter{Email: ems})
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

func (api *UseCase) GetAllInterests(ctx context.Context) ([]*auth.Interest, error) {
	return api.interestStorage.Get(ctx, nil)
}

func (api *UseCase) Registration(body auth.RegitstrationBody, ctx context.Context) (string, types.UserID, error) {
	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		return "", -1, err
	}

	err = api.dbReader.AddAccount(ctx, body.Name, body.Birthday, body.Gender, body.Email, hashedPassword)
	if err != nil {
		return "", -1, err
	}

	SID, UID, err := api.Login(body.Email, body.Password, ctx)
	if err != nil {
		return "", -1, err
	}

	interests, err := api.interestStorage.Get(ctx, &auth.InterestGetFilter{Name: body.Interests})
	if err != nil {
		return "", -1, err
	}
	err = api.interestStorage.CreatePersonInterests(ctx, UID, getInterestIDs(interests))
	if err != nil {
		return SID, UID, err
	}
	return SID, UID, nil
}

func (api *UseCase) GetName(sessionID string, ctx context.Context) (string, error) {
	person, err := api.dbReader.Get(ctx, &auth.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		return "", err
	}

	if len(person) == 0 {
		return "", errors.New("no person with such sessionID")
	}

	return person[0].Name, err
}

func getInterestIDs(interests []*auth.Interest) []types.InterestID {
	res := make([]types.InterestID, len(interests))
	for i := range interests {
		res[i] = interests[i].ID
	}
	return res
}

func (api *UseCase) Logout(sessionID string, ctx context.Context) error {
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
