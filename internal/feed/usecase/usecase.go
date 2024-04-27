package usecase

import (
	"context"

	"main.go/internal/feed"
	"main.go/internal/types"
)

type UseCase struct {
	personStorage   feed.PersonStorage
	interestStorage feed.InterestStorage
	imageStorage    feed.ImageStorage
	likeStorage     feed.LikeStorage
}

func New(pstor feed.PersonStorage, istor feed.InterestStorage, imstor feed.ImageStorage, lstor feed.LikeStorage) *UseCase {
	return &UseCase{
		personStorage:   pstor,
		interestStorage: istor,
		imageStorage:    imstor,
		likeStorage:     lstor,
	}
}

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *UseCase) GetCards(userID types.UserID, ctx context.Context) ([]feed.Card, error) {
	persons, err := service.personStorage.GetFeed(ctx, userID)

	if err != nil {
		return nil, err
	}

	interests, images, err := service.getUserCards(persons, ctx)

	if err != nil {
		return nil, err
	}

	return combineToCards(persons, interests, images), nil
}

func (service *UseCase) CreateLike(profile1, profile2 types.UserID, ctx context.Context) (int, error) {
	// likes, err := service.likeStorage.GetLikesLeft(ctx, profile1)
	// if err != nil {
	// 	return 0, err
	// }
	likes, err := service.likeStorage.DecreaseLikesCount(ctx, profile1)
	if err != nil {
		return 0, err
	}
	if likes == 0 {
		err := service.likeStorage.IncreaseLikesCount(ctx, profile1)
		if err != nil {
			return 0, err
		}
	}
	return likes, service.likeStorage.Create(ctx, profile1, profile2)
}

func (service *UseCase) getUserCards(persons []*feed.Person, ctx context.Context) ([][]*feed.Interest, [][]feed.Image, error) {
	var err error
	interests := make([][]*feed.Interest, len(persons))
	images := make([][]feed.Image, len(persons))
	for j := range persons {
		interests[j], err = service.interestStorage.GetPersonInterests(ctx, persons[j].ID)
		if err != nil {
			return nil, nil, err
		}
		images[j], err = service.imageStorage.Get(ctx, int64(persons[j].ID))
		if err != nil {
			return nil, nil, err
		}
	}
	return interests, images, nil
}

func combineToCards(persons []*feed.Person, interests [][]*feed.Interest, images [][]feed.Image) []feed.Card {
	if len(persons) != len(interests) || len(persons) != len(images) {
		return nil
	}

	photos := make([][]feed.ImageToSend, len(persons))
	for i := range images {
		photos[i] = make([]feed.ImageToSend, len(images[i]))
		for j, image := range images[i] {
			photos[i][j] = feed.ImageToSend{
				Cell: image.CellNumber,
				Url:  image.Url,
			}
		}
	}

	res := make([]feed.Card, len(persons))
	for i := range persons {
		res[i] = feed.Card{ID: persons[i].ID, Name: persons[i].Name, Birthday: persons[i].Birthday, Description: persons[i].Description,
			Interests: interests[i], Photos: photos[i]}
	}
	return res
}
