package delivery

import (
	"context"
	models "main.go/db"
	requests "main.go/internal/pkg"
	"main.go/internal/service"
	"main.go/internal/types"
)

type (
	Service interface {
		GetCards(userID types.UserID, ctx context.Context) ([]models.Card, error) //
		GetName(sessionID string, ctx context.Context) (string, error)            //
		GetAllInterests(ctx context.Context) ([]*models.Interest, error)          //
		GetProfile(params service.ProfileGetParams, ctx context.Context) ([]models.Card, error)
		UpdateProfile(SID string, profile requests.ProfileUpdateRequest, ctx context.Context) error //
		DeleteProfile(sessionID string, ctx context.Context) error                                  //
		CreateLike(profile1, profile2 types.UserID, ctx context.Context) error                      //
		GetMatches(profile types.UserID, ctx context.Context) ([]models.Card, error)
		GetImage(userID int64, ctx context.Context) ([]models.Image, error) //
		AddImage(userImage models.Image, ctx context.Context) error         //
		DeleteImage(userImage models.Image, ctx context.Context) error      //
	}

	Auth interface {
		IsAuthenticated(sessionID string, ctx context.Context) (types.UserID, bool)
		Login(email, password string, ctx context.Context) (string, types.UserID, error)
		Logout(sessionID string, ctx context.Context) error
		Registration(Name string, Birthday string, Gender string, Email string, Password string, ctx context.Context) (string, types.UserID, error)
	}
)
