package auth

import (
	"context"

	"main.go/internal/types"
)

type (
	IUseCase interface {
		IsAuthenticated(sessionID string, ctx context.Context) (types.UserID, bool, error)
		Login(email, password string, ctx context.Context) (*Profile, string, error)
		Logout(sessionID string, ctx context.Context) error
		Registration(body RegitstrationBody, ctx context.Context) (*Profile, string, error)
		GetAllInterests(ctx context.Context) ([]*Interest, error)
		GetName(sessionID string, ctx context.Context) (string, error)
		GetProfile(params ProfileGetParams, ctx context.Context) ([]Profile, error)
		UpdateProfile(SID string, profile ProfileUpdateRequest, ctx context.Context) error
		DeleteProfile(sessionID string, ctx context.Context) error
		GetMatches(profile types.UserID, ctx context.Context) ([]Profile, error)
	}
	PersonStorage interface {
		AddAccount(ctx context.Context, Name string, Birhday string, Gender string, Email string, Password string) error
		Get(ctx context.Context, filter *PersonGetFilter) ([]*Person, error)
		Update(ctx context.Context, person Person) error
		Delete(ctx context.Context, sessionID string) error
		RemoveSession(ctx context.Context, sid string) error
		GetMatch(ctx context.Context, person1ID types.UserID) ([]types.UserID, error)
	}

	InterestStorage interface {
		Get(ctx context.Context, filter *InterestGetFilter) ([]*Interest, error)
		CreatePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error
		GetPersonInterests(ctx context.Context, personID types.UserID) ([]*Interest, error)
		DeletePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error
	}

	ImageStorage interface {
		Get(ctx context.Context, userID int64) ([]Image, error)
	}
)
