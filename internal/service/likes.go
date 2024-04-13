package service

import (
	"context"
	models "main.go/db"
	"main.go/internal/types"
)

func (service *Service) CreateLike(profile1, profile2 types.UserID, ctx context.Context) error {
	return service.likeStorage.Create(ctx, profile1, profile2)
}

func (service *Service) GetMatches(profile types.UserID, ctx context.Context) ([]models.Card, error) {
	ids, err := service.likeStorage.GetMatch(ctx, profile)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	return service.GetProfile(ProfileGetParams{ID: ids}, ctx)
}
