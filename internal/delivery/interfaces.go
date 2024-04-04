package delivery

import models "main.go/db"

type Service interface {
	GetCards(sessionID string, reqID int64) (string, error)
	GetName(sessionID string, reqID int64) (string, error)
	GetAllInterests(reqID int64) (string, error)
	GetImage(sessionID string, requestID int64) (string, error)
	AddImage(userImage models.Image, requestID int64) error
	GetId(sessionID string, requestID int64) (int, error)
}

type Auth interface {
	IsAuthenticated(sessionID string, reqID int64) bool
	Login(email, password string, reqID int64) (string, error)
	Logout(sessionID string, reqID int64) error
	Registration(Name string, Birthday string, Gender string, Email string, Password string, reqID int64) (string, error)
}
