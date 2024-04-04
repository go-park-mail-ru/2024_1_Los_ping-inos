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
	imageStorage    ImageStorage
}

func New(pstor PersonStorage, istor InterestStorage, imstor ImageStorage) *Service {
	return &Service{
		personStorage:   pstor,
		interestStorage: istor,
		imageStorage:    imstor,
	}
}

func (service *Service) GetName(sessionID string, requestID int64) (string, error) {
	person, err := service.personStorage.Get(requestID, &models.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		return "", err
	}

	if person == nil || len(person) == 0 {
		return "", errors.New("no person with such sessionID")
	}

	return person[0].Name, err
}

func (service *Service) GetId(sessionID string, requestID int64) (int, error) {
	person, err := service.personStorage.Get(requestID, &models.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		return 0, err
	}

	if person == nil || len(person) == 0 {
		return 0, errors.New("no person with such sessionID")
	}

	return int(person[0].ID), err
}

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *Service) GetCards(sessionID string, requestID int64) (string, error) {
	persons, err := service.personStorage.Get(requestID, nil)

	if err != nil {
		return "", err
	}

	i := 0
	for ; i < len(persons); i++ { // :eyes:
		if persons[i].SessionID == sessionID {
			persons = append(persons[:i], persons[i+1:]...)
			break
		}
	}

	res, _ := personsToJSON(persons)

	return res, nil
}

func personsToJSON(persons []*models.Person) (string, error) {
	res, err := json.Marshal(persons)
	return string(res), err
}
