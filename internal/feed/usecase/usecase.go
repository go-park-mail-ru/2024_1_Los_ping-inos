package usecase

import (
	"context"
	"main.go/internal/feed"
	"main.go/internal/types"
)

type UseCase struct {
	storage feed.PostgresStorage
	ws      feed.WebSocStorage
}

func New(pstor feed.PostgresStorage, wsstor feed.WebSocStorage) *UseCase {
	return &UseCase{
		storage: pstor,
		ws:      wsstor,
	}
}

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *UseCase) GetCards(userID types.UserID, ctx context.Context) ([]feed.Card, error) {
	persons, err := service.storage.GetFeed(ctx, userID)

	if err != nil {
		return nil, err
	}

	interests, images, err := service.getUserCards(persons, ctx)

	if err != nil {
		return nil, err
	}

	return combineToCards(persons, interests, images), nil
}

func (service *UseCase) CreateLike(profile1, profile2 types.UserID, ctx context.Context) error {
	return service.storage.CreateLike(ctx, profile1, profile2)
}

func (service *UseCase) GetChat(ctx context.Context, user1, user2 types.UserID) ([]feed.Message, error) {
	return service.storage.GetChat(ctx, user1, user2)
}

func (service *UseCase) SaveMessage(ctx context.Context, message feed.Message) (*feed.Message, error) {
	return service.storage.CreateMessage(ctx, message)
}

func (service *UseCase) getUserCards(persons []*feed.Person, ctx context.Context) ([][]*feed.Interest, [][]feed.Image, error) {
	var err error
	interests := make([][]*feed.Interest, len(persons))
	images := make([][]feed.Image, len(persons))
	for j := range persons {
		interests[j], err = service.storage.GetPersonInterests(ctx, persons[j].ID)
		if err != nil {
			return nil, nil, err
		}
		images[j], err = service.storage.GetImages(ctx, int64(persons[j].ID))
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

func (service *UseCase) CreateClaim(ctx context.Context, typeID, senderID, receiverID int64) error {
	return service.storage.CreateClaim(ctx, feed.Claim{TypeID: typeID, SenderID: senderID, ReceiverID: receiverID})
}