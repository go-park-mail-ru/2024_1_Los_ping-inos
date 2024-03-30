package service

import (
	"errors"
	"github.com/sirupsen/logrus"
	models "main.go/db"
	. "main.go/internal/logs"
	"time"
)

func (service *Service) GetProfile(sessionID string, requestID int64) (string, error) {
	person, err := service.personStorage.Get(requestID, &models.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		return "", err
	}

	if len(person) == 0 {
		return "", errors.New("no person with such sessionID")
	}

	interests, err := service.interestStorage.GetPersonInterests(requestID, person[0].ID)
	if err != nil {
		return "", err
	}

	res, err := personsToJSON(person, [][]*models.Interest{interests})
	if err != nil {
		return "", err
	}
	return res, err
}

func (service *Service) UpdateProfile(sessionID, name, password, description, birthday string, interests []string, requestID int64) error {
	persons, err := service.personStorage.Get(requestID, &models.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		return err
	}
	person := persons[0]
	if name != "" {
		person.Name = name
	}
	if birthday != "" {
		person.Birthday, err = time.Parse("01.02.2006", birthday)
	}
	if description != "" {
		person.Description = description
	}
	if password != "" {
		person.Password, err = hashPassword(password)
		if err != nil {
			return err
		}
	}
	if interests != nil {
		// TODO
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("profile update interest not implemented")
	}
	err = service.personStorage.Update(requestID, *person)

	return err
}

func (service *Service) DeleteProfile(sessionID string, requestID int64) error {
	err := service.personStorage.Delete(requestID, sessionID)
	return err
}
