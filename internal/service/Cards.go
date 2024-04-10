package service

import (
	"encoding/json"
	"errors"
	"main.go/internal/types"

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

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *Service) GetCards(userID types.UserID, requestID int64) (string, error) {
	persons, err := service.personStorage.GetFeed(requestID, userID)

	if err != nil {
		return "", err
	}

	interests, images, err := service.getUserCards(persons, requestID)
	if err != nil {
		return "", err
	}

	return personsToJSON(persons, interests, images)
}

func (service *Service) getUserCards(persons []*models.Person, requestID int64) ([][]*models.Interest, [][]string, error) {
	var err error
	interests := make([][]*models.Interest, len(persons))
	images := make([][]string, len(persons))
	for j := range persons {
		interests[j], err = service.interestStorage.GetPersonInterests(requestID, persons[j].ID)
		if err != nil {
			return nil, nil, err
		}
		tmp, err := service.imageStorage.Get(requestID, int64(persons[j].ID))
		if err != nil {
			return nil, nil, err
		}
		for t := range tmp {
			images[j] = append(images[j], tmp[t].Url)
		}
	}
	return interests, images, nil
}

func combineToCards(persons []*models.Person, interests [][]*models.Interest, images [][]string) []models.Card {
	if len(persons) != len(interests) || len(persons) != len(images) {
		Log.Warn("can't create cards: different slices size")
		return nil
	}
	res := make([]models.Card, len(persons))
	for i := range persons {
		res[i] = models.Card{Person: persons[i], Interests: interests[i], Photo: images[i]}
	}
	return res
}

func personsToJSON(persons []*models.Person, interests [][]*models.Interest, images [][]string) (string, error) {
	combined := combineToCards(persons, interests, images)
	if combined == nil {
		return "", errors.New("can't create cards: different persons and interests sizes")
	}
	res, err := json.Marshal(combined)

	return string(res), err
}
