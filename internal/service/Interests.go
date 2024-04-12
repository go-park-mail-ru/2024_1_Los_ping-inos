package service

import (
	models "main.go/db"
)

func (service *Service) GetAllInterests(requestID int64) ([]*models.Interest, error) {
	return service.interestStorage.Get(requestID, nil)
}
