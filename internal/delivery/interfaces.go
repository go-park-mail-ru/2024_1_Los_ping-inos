package delivery

import (
	requests "main.go/internal/pkg"
	"main.go/internal/service"
)

type (
	Service interface {
		GetCards(sessionID string, reqID int64) (string, error)
		GetName(sessionID string, reqID int64) (string, error)
		GetAllInterests(reqID int64) (string, error)
		GetProfile(params service.ProfileGetParams, requestID int64) (string, error)
		UpdateProfile(SID string, profile requests.ProfileUpdateRequest, requestID int64) error
		DeleteProfile(sessionID string, requestID int64) error
	}

	Auth interface {
		IsAuthenticated(sessionID string, reqID int64) bool
		Login(email, password string, reqID int64) (string, error)
		Logout(sessionID string, reqID int64) error
		Registration(Name string, Birthday string, Gender string, Email string, Password string, reqID int64) (string, error)
	}
)
