package feed

import (
	"context"
	"main.go/internal/types"
)

type (
	UseCase interface {
		GetCards(userID types.UserID, ctx context.Context) ([]Card, error)
		CreateLike(profile1, profile2 types.UserID, ctx context.Context) error
	}
	PostgresStorage interface {
		GetFeed(ctx context.Context, filter types.UserID) ([]*Person, error)
		GetPersonInterests(ctx context.Context, personID types.UserID) ([]*Interest, error)
		GetLike(ctx context.Context, filter *LikeGetFilter) ([]types.UserID, error)
		CreateLike(ctx context.Context, person1ID, person2ID types.UserID) error
		GetImages(ctx context.Context, userID int64) ([]Image, error)
	}
)
