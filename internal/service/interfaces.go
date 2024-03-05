package service

import models "main.go/db"

type PersonStorage interface {
	Get(filter *models.PersonGetFilter) ([]*models.Person, error)
	AddAccount(Name string, Birhday string, Gender string, Email string, Password string) error
	Update(person models.Person) error
}

type InterestStorage interface {
	Get() ([]*models.Interest, error)
}
