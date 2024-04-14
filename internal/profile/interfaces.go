package profile

import (
	"context"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
)

type (
	UseCase interface {
		GetProfile(params ProfileGetParams, ctx context.Context) ([]Card, error)
		UpdateProfile(SID string, profile requests.ProfileUpdateRequest, ctx context.Context) error
		DeleteProfile(sessionID string, ctx context.Context) error
		GetMatches(profile types.UserID, ctx context.Context) ([]Card, error)
	}

	PersonStorage interface {
		Get(ctx context.Context, filter *PersonGetFilter) ([]*Person, error)
		Update(ctx context.Context, person Person) error
		Delete(ctx context.Context, sessionID string) error
	}

	InterestStorage interface {
		Get(ctx context.Context, filter *InterestGetFilter) ([]*Interest, error)
		GetPersonInterests(ctx context.Context, personID types.UserID) ([]*Interest, error)
		CreatePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error
		DeletePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error
	}
	ImageStorage interface {
		Get(ctx context.Context, userID int64) ([]Image, error)
	}
	LikeStorage interface {
		Get(ctx context.Context, filter *LikeGetFilter) ([]types.UserID, error)
		Create(ctx context.Context, person1ID, person2ID types.UserID) error
		GetMatch(ctx context.Context, person1ID types.UserID) ([]types.UserID, error)
	}
)
