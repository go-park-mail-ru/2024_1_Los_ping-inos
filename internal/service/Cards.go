package service

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
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

func (service *Service) GetName(sessionID string) (string, error) {
	person, err := service.personStorage.Get(&models.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		logrus.Info(err.Error())
		return "", err
	}
	logrus.Info(len(person))
	if person == nil || len(person) == 0 {
		logrus.Info("No person with such session ID ", sessionID)
		return "", err
	}

	return person[0].Name, err
}

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *Service) GetCards(sessionID string) (string, error) {
	persons, err := service.personStorage.Get(nil)

	if err != nil {
		return "", errors.New("can't get users")
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
