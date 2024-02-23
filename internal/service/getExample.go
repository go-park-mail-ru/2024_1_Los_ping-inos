package service

import (
	"main.go/internal/types"
)

// Storage - С ростом количества крудов тоже растёт
type Storage interface {
	Get(id types.UserID) string
}

// Service - Обработчик всей логики
type Service struct {
	storage Storage
}

func New(s Storage) *Service {
	return &Service{storage: s}
}

// GetCoolIdsList - пример логики
func (e *Service) GetCoolIdsList() ([]string, error) {
	res := make([]string, 0)
	var i types.UserID
	for ; i < 5; i++ {
		res = append(res, e.storage.Get(i))
	}
	return res, nil
}
