package service

import (
	"errors"
	"main.go/db"
	"main.go/internal/types"
	"net/http"
)

// Storage - С ростом количества крудов тоже растёт
type Storage interface {
	PersonStorage
}

type PersonStorage interface {
	Get(filter *models.PersonFilter) ([]*models.Person, error)
}

// Service - Обработчик всей логики
type Service struct {
	storage Storage
	auth    Auth
}

func New(stor Storage, auth Auth) *Service {
	return &Service{storage: stor, auth: auth}
}

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *Service) GetCards(w http.ResponseWriter, r *http.Request) ([]models.Person, error) {
	if !service.auth.IsAuthenticated(w, r) {
		return nil, errors.New("not authenticated")
	}

	res := make([]models.Person, 0)
	var i types.UserID
	for ; i < 5; i++ {
		// TODO логику
	}
	return res, nil
}
