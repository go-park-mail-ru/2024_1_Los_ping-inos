package service

import (
	"main.go/db"
	"main.go/internal/types"
)

type PersonStorage interface {
	Get(filter *models.PersonFilter) ([]*models.Person, error)
}

// Service - Обработчик всей логики
type Service struct {
	storage PersonStorage
}

func New(stor PersonStorage) *Service {
	return &Service{storage: stor}
}

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *Service) GetCards(sessionID string) ([]models.Person, error) {
	res := make([]models.Person, 0)
	var i types.UserID
	for ; i < 5; i++ {
		// TODO логику
	}
	return res, nil
}
