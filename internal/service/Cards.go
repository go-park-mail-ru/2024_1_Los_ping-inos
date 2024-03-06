package service

import (
	"encoding/json"
	"errors"
	models "main.go/db"
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
func (service *Service) GetCards(sessionID string) (string, error) {

	persons, err := service.personStorage.Get(nil)

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
