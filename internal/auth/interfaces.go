package auth

import (
	"context"

	"google.golang.org/grpc"
	image "main.go/internal/image/protos/gen"
	"main.go/internal/types"
)

type (
	IUseCase interface {
		IsAuthenticated(sessionID string, ctx context.Context) (types.UserID, bool, error)
		Login(email, password string, ctx context.Context) (*Profile, string, error)
		Logout(sessionID string, ctx context.Context) error
		Registration(body RegitstrationBody, ctx context.Context) (*Profile, string, error)
		GetAllInterests(ctx context.Context) ([]*Interest, error)
		GetName(userID types.UserID, ctx context.Context) (string, error)
		GetProfile(params ProfileGetParams, ctx context.Context) ([]Profile, error)
		UpdateProfile(UID types.UserID, profile ProfileUpdateRequest, ctx context.Context) error
		DeleteProfile(UID types.UserID, ctx context.Context) error
		GetMatches(profile types.UserID, nameFilter string, ctx context.Context) ([]Profile, error)
	}
	PersonStorage interface {
		AddAccount(ctx context.Context, Name string, Birthday string, Gender string, Email string, Password string) (string, error)
		Get(ctx context.Context, filter *PersonGetFilter) ([]*Person, error)
		Update(ctx context.Context, person Person) error
		Delete(ctx context.Context, UID types.UserID) error
		GetMatch(ctx context.Context, person1ID types.UserID) ([]types.UserID, error)
		//GetUserCards(ctx context.Context, persons []types.UserID) ([][]*Interest, [][]Image, error)
	}

	SessionStorage interface {
		GetBySID(ctx context.Context, SID string) (*Session, error)
		CreateSession(ctx context.Context, UID types.UserID) (string, error)
		DeleteSession(ctx context.Context, SID string) error
	}

	InterestStorage interface {
		Get(ctx context.Context, filter *InterestGetFilter) ([]*Interest, error)
		CreatePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error
		GetPersonInterests(ctx context.Context, personID types.UserID) ([]*Interest, error)
		DeletePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error
	}

	ImageClient interface {
		GetImage(ctx context.Context, in *image.GetImageRequest, opts ...grpc.CallOption) (*image.GetImageResponce, error)
	}
)
