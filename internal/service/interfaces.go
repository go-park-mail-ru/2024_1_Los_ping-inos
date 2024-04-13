package service

import (
	"context"
	models "main.go/db"
	"main.go/internal/types"
)

type PersonStorage interface {
	Get(ctx context.Context, filter *models.PersonGetFilter) ([]*models.Person, error)
	GetFeed(ctx context.Context, filter types.UserID) ([]*models.Person, error)
	AddAccount(ctx context.Context, Name string, Birhday string, Gender string, Email string, Password string) error
	Update(ctx context.Context, person models.Person) error
	RemoveSession(ctx context.Context, sid string) error
	Delete(ctx context.Context, sessionID string) error
}

type InterestStorage interface {
	Get(ctx context.Context, filter *models.InterestGetFilter) ([]*models.Interest, error)
	GetPersonInterests(ctx context.Context, personID types.UserID) ([]*models.Interest, error)
	CreatePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error
	DeletePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error
}

type LikeStorage interface {
	Get(ctx context.Context, filter *models.LikeGetFilter) ([]types.UserID, error)
	Create(ctx context.Context, person1ID, person2ID types.UserID) error
	GetMatch(ctx context.Context, person1ID types.UserID) ([]types.UserID, error)
}

type ImageStorage interface {
	Get(ctx context.Context, userID int64) ([]models.Image, error)
	Add(ctx context.Context, image models.Image) error
	Delete(ctx context.Context, image models.Image) error
}
