package service

import (
	"errors"
	"fmt"
	models "main.go/db"
	. "main.go/internal/logs"
	"main.go/internal/types"
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

	if len(person) == 0 {
		return "", errors.New("no person with such sessionID")
	}

	return person[0].Name, err
}

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *Service) GetCards(userID types.UserID, requestID int64) ([]models.Card, error) {
	persons, err := service.personStorage.GetFeed(requestID, userID)

	if err != nil {
		return nil, err
	}

	interests, images, err := service.getUserCards(persons, requestID)

	if err != nil {
		return nil, err
	}

	return combineToCards(persons, interests, images), nil
}

func (service *Service) getUserCards(persons []*models.Person, requestID int64) ([][]*models.Interest, [][]models.Image, error) {
	var err error
	interests := make([][]*models.Interest, len(persons))
	images := make([][]models.Image, len(persons))
	for j := range persons {
		interests[j], err = service.interestStorage.GetPersonInterests(requestID, persons[j].ID)
		if err != nil {
			return nil, nil, err
		}
		images[j], err = service.imageStorage.Get(requestID, int64(persons[j].ID))
		fmt.Printf("%v", images[j])
		if err != nil {
			return nil, nil, err
		}
	}
	return interests, images, nil
}

func combineToCards(persons []*models.Person, interests [][]*models.Interest, images [][]models.Image) []models.Card {
	if len(persons) != len(interests) || len(persons) != len(images) {
		Log.Warn("can't create cards: different slices size")
		return nil
	}

	photos := make([][]models.ImageToSend, len(persons))
	for i := range images {
		photos[i] = make([]models.ImageToSend, len(images[i]))
		for j, image := range images[i] {
			photos[i][j] = models.ImageToSend{
				Cell: image.CellNumber,
				Url:  image.Url,
			}
		}
	}

	res := make([]models.Card, len(persons))
	for i := range persons {
		res[i] = models.Card{Name: persons[i].Name, Birthday: persons[i].Birthday, Description: persons[i].Description,
			Interests: interests[i], Photos: photos[i]}
	}
	return res
}
