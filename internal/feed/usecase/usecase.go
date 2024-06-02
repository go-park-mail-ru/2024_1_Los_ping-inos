package usecase

import (
	"context"
	"fmt"

	"main.go/internal/feed"
	image "main.go/internal/image/protos/gen"
	"main.go/internal/types"
)

type UseCase struct {
	storage    feed.PostgresStorage
	ws         feed.WebSocStorage
	grpcClient image.ImageClient
}

func New(pstor feed.PostgresStorage, wsstor feed.WebSocStorage, grpcClient image.ImageClient) *UseCase {
	return &UseCase{
		storage:    pstor,
		ws:         wsstor,
		grpcClient: grpcClient,
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
	likesLeft, err := service.storage.DecreaseLikesCount(ctx, profile1)
	if err != nil {
		return err
	}
	if likesLeft >= 0 {
		return service.storage.CreateLike(ctx, profile1, profile2)
	}
	return feed.NoLikesLeftErr
}

func (service *UseCase) GetChat(ctx context.Context, user1, user2 types.UserID) ([]feed.Message, []feed.Image, []feed.Person, error) {
	imagePerson := []feed.Image{}
	for i := 0; i < 5; i++ {
		image, err := service.grpcClient.GetImage(ctx, &image.GetImageRequest{Id: int64(user1), Cell: fmt.Sprintf("%v", i)})
		imagePiece := feed.Image{}
		if err != nil {
			imagePiece = feed.Image{
				UserId:     int64(user1),
				Url:        "",
				CellNumber: fmt.Sprintf("%v", i),
			}
		} else {
			imagePiece = feed.Image{
				UserId:     int64(user1),
				Url:        image.Url,
				CellNumber: fmt.Sprintf("%v", i),
			}
		}
		imagePerson = append(imagePerson, imagePiece)
	}

	persons, err := service.storage.GetPerson(ctx, user1)
	if err != nil {
		return nil, nil, nil, err
	}
	msgs, err := service.storage.GetChat(ctx, user1, user2)
	if err != nil {
		return nil, nil, nil, err
	}

	return msgs, imagePerson, persons, nil
}

func (service *UseCase) SaveMessage(ctx context.Context, message feed.MessageToReceive) (*feed.MessageToReceive, error) {
	return service.storage.CreateMessage(ctx, message)
}

func (service *UseCase) getUserCards(persons []*feed.Person, ctx context.Context) ([][]*feed.Interest, [][]feed.Image, error) {
	var err error
	interests := make([][]*feed.Interest, len(persons))
	images := make([][]feed.Image, len(persons))

	for j := range persons {
		imagePerson := []feed.Image{}
		for i := 0; i < 5; i++ {
			image, err := service.grpcClient.GetImage(ctx, &image.GetImageRequest{Id: int64(persons[j].ID), Cell: fmt.Sprintf("%v", i)})
			imagePiece := feed.Image{}
			if err != nil {
				imagePiece = feed.Image{
					UserId:     int64(persons[j].ID),
					Url:        "",
					CellNumber: fmt.Sprintf("%v", i),
				}
			} else {
				imagePiece = feed.Image{
					UserId:     int64(persons[j].ID),
					Url:        image.Url,
					CellNumber: fmt.Sprintf("%v", i),
				}
			}
			imagePerson = append(imagePerson, imagePiece)
		}
		images[j] = imagePerson

		interests[j], err = service.storage.GetPersonInterests(ctx, persons[j].ID)
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

func (service *UseCase) GetClaims(ctx context.Context) ([]feed.PureClaim, error) {
	return service.storage.GetAllClaims(ctx)
}
