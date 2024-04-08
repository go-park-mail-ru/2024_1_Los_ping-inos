package service

import (
	"encoding/json"
	"errors"

	models "main.go/db"
	. "main.go/internal/logs"
)

// Service - Обработчик всей логики
type Service struct {
	personStorage   PersonStorage
	interestStorage InterestStorage
	imageStorage    ImageStorage
  likeStorage     LikeStorage
}

func New(pstor PersonStorage, istor InterestStorage, imstor ImageStorage, lstor LikeStorage) *Service {
	return &Service{
		personStorage:   pstor,
		interestStorage: istor,
		imageStorage:    imstor,
    likeStorage:     lstor,
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

func (service *Service) GetId(sessionID string, requestID int64) (int64, error) {
	person, err := service.personStorage.Get(requestID, &models.PersonGetFilter{SessionID: []string{sessionID}})
	if err != nil {
		return 0, err
	}

	if person == nil || len(person) == 0 {
		return 0, errors.New("no person with such sessionID")
	}

	return int64(person[0].ID), err
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

	interests := make([][]*models.Interest, len(persons))
	for j := range persons {
		interests[j], err = service.interestStorage.GetPersonInterests(requestID, persons[i].ID)
		if err != nil {
			return "", err
		}
	}

	return personsToJSON(persons, interests)
}

func combinePersonsAndInterestsToCards(persons []*models.Person, interests [][]*models.Interest) []models.PersonWithInterests {
	if len(persons) != len(interests) {
		Log.Warn("can't create cards: different slices size")
		return nil
	}
	res := make([]models.PersonWithInterests, len(persons))
	for i := range persons {
		res[i] = models.PersonWithInterests{Person: persons[i], Interests: interests[i]}
	}
	return res
}

func personsToJSON(persons []*models.Person, interests [][]*models.Interest) (string, error) {
	combined := combinePersonsAndInterestsToCards(persons, interests)
	if combined == nil {
		return "", errors.New("can't create cards: different persons and interests sizes")
	}
	res, err := json.Marshal(combined)

	return string(res), err
}
