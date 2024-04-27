package feed

import (
	"context"

	"main.go/internal/types"
)

type (
	UseCase interface {
		GetCards(userID types.UserID, ctx context.Context) ([]Card, error)
		CreateLike(profile1, profile2 types.UserID, ctx context.Context) (int, error)
	}

	PersonStorage interface {
		GetFeed(ctx context.Context, filter types.UserID) ([]*Person, error)
	}

	InterestStorage interface {
		GetPersonInterests(ctx context.Context, personID types.UserID) ([]*Interest, error)
	}

	LikeStorage interface {
		Get(ctx context.Context, filter *LikeGetFilter) ([]types.UserID, error)
		Create(ctx context.Context, person1ID, person2ID types.UserID) error
		GetLikesLeft(ctx context.Context, personID types.UserID) (int, error)
		DecreaseLikesCount(ctx context.Context, personID types.UserID) (int, error)
		IncreaseLikesCount(ctx context.Context, personID types.UserID) error
	}

	ImageStorage interface {
		Get(ctx context.Context, userID int64) ([]Image, error)
	}
)
