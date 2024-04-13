package service

import (
	"context"
	models "main.go/db"
)

func (service *Service) GetAllInterests(ctx context.Context) ([]*models.Interest, error) {
	return service.interestStorage.Get(ctx, nil)
}
