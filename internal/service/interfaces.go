package service

import (
	models "main.go/db"
	"main.go/internal/types"
)

type PersonStorage interface {
	Get(requestID int64, filter *models.PersonGetFilter) ([]*models.Person, error)
	AddAccount(requestID int64, Name string, Birhday string, Gender string, Email string, Password string) error
	Update(requestID int64, person models.Person) error
	RemoveSession(requestID int64, sid string) error
	Delete(requestID int64, sessionID string) error
}

type InterestStorage interface {
	Get(requestID int64, filter *models.InterestGetFilter) ([]*models.Interest, error)
	GetPersonInterests(requestID int64, personID types.UserID) ([]*models.Interest, error)
	CreatePersonInterests(requestID int64, personID types.UserID, interestID []types.InterestID) error
	DeletePersonInterests(requestID int64, personID types.UserID, interestID []types.InterestID) error
}

type LikeStorage interface {
	Get(requestID int64, filter *models.LikeGetFilter) ([]*models.Like, error)
	Create(requestID int64, person1ID, person2ID types.UserID) error
	GetMatch(requestID int64, person1ID types.UserID) ([]types.UserID, error)
}

type ImageStorage interface {
	Get(requestID int64, userID int64) ([]models.Image, error)
	Add(requestID int64, image models.Image) error
	Delete(requestID int64, image models.Image) error
}
