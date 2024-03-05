package service

import (
	"encoding/json"
	"errors"
	"golang.org/x/exp/slices"
	models "main.go/db"
	"main.go/internal/types"
)

type PersonStorage interface {
	Get(filter *models.PersonGetFilter) ([]*models.Person, error)
	AddAccount(Name string, Birhday string, Gender string, Email string, Password string) error
	Update(person models.Person) error
}

// Service - Обработчик всей логики
type Service struct {
	storage PersonStorage
}

func New(stor PersonStorage) *Service {
	return &Service{storage: stor}
}

// GetCards - вернуть ленту пользователей, доступно только авторизованному пользователю
func (service *Service) GetCards(sessionID string, firstID types.UserID) (string, error) {
	ids := make([]types.UserID, 5) // пока что без какой-либо логики возвращаем первые 5 от последнего запроса
	for i := range ids {
		ids[i] = firstID + types.UserID(i+1)
	}

	user, err := service.storage.Get(&models.PersonGetFilter{SessionID: []string{sessionID}})

	if (err != nil || user == nil) && slices.Contains(ids, user[0].ID) {
		ids[slices.Index(ids, user[0].ID)] = firstID + 6
	}

	persons, err := service.storage.Get(&models.PersonGetFilter{ID: ids})

	if err != nil {
		return "", errors.New("can't get users")
	}

	res, _ := personsToJSON(persons)

	return res, nil
}

func personsToJSON(persons []*models.Person) (string, error) {
	res, err := json.Marshal(persons)
	return string(res), err
}
