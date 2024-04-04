package service

import models "main.go/db"

type PersonStorage interface {
	Get(requestID int64, filter *models.PersonGetFilter) ([]*models.Person, error)
	AddAccount(requestID int64, Name string, Birhday string, Gender string, Email string, Password string) error
	Update(requestID int64, person models.Person) error
	RemoveSession(requestID int64, sid string) error
}

type InterestStorage interface {
	Get(requestID int64) ([]*models.Interest, error)
}

type ImageStorage interface {
	Get(requestID int64, person models.Person) (*models.Image, error)
	Add(requestID int64, image models.Image) error
}
