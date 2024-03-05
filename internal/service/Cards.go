package service

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	models "main.go/db"
	"main.go/internal/types"
)

// Service - Обработчик всей логики
type Service struct {
	personStorage   PersonStorage
	interestStorage InterestStorage
}

func New(pstor PersonStorage, istor InterestStorage) *Service {
	return &Service{
		personStorage:   pstor,
		interestStorage: istor,
	}
}

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *Service) GetCards(sessionID string, firstID types.UserID) (string, error) {
	ids := make([]types.UserID, 5) // пока что без какой-либо логики возвращаем первые 5 от последнего запроса
	for i := range ids {
		ids[i] = firstID + types.UserID(i+1)
	}

	user, err := service.personStorage.Get(&models.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		logrus.Info(err.Error())
		return "", err
	}

	var ID types.UserID
	if user == nil {
		ID = 0
	} else {
		ID = user[0].ID
	}

	if slices.Contains(ids, ID) {
		ids[slices.Index(ids, ID)] = firstID + 6
	}

	persons, err := service.personStorage.Get(&models.PersonGetFilter{ID: ids})

	if err != nil {
		return "", errors.New("can't get users")
	}

	res, _ := personsToJSON(persons)

	return res, nil
}

func personsToJSON(persons []*models.Person) (string, error) {
	res, err := json.Marshal(persons)
	return string(res), err
}
