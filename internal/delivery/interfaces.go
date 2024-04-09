package delivery

import (
	models "main.go/db"
	requests "main.go/internal/pkg"
	"main.go/internal/service"
	"main.go/internal/types"
)

type (
	Service interface {
		GetCards(sessionID string, reqID int64) (string, error)
		GetName(sessionID string, reqID int64) (string, error)
		GetAllInterests(reqID int64) (string, error)
		GetProfile(params service.ProfileGetParams, requestID int64) (string, error)
		UpdateProfile(SID string, profile requests.ProfileUpdateRequest, requestID int64) error
		DeleteProfile(sessionID string, requestID int64) error
		CreateLike(profile1, profile2 types.UserID, requestID int64) error
		GetMatches(profile types.UserID, requestID int64) (string, error)
		GetImage(userID int64, requestID int64) ([]models.Image, error)
		AddImage(userImage models.Image, requestID int64) error
		DeleteImage(userImage models.Image, requestID int64) error
	}

	Auth interface {
		IsAuthenticated(sessionID string, reqID int64) (types.UserID, bool)
		Login(email, password string, reqID int64) (string, error)
		Logout(sessionID string, reqID int64) error
		Registration(Name string, Birthday string, Gender string, Email string, Password string, reqID int64) (string, error)
	}
)
