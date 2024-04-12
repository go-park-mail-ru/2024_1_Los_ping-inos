package service

import (
	models "main.go/db"
	"main.go/internal/types"
)

func (service *Service) CreateLike(profile1, profile2 types.UserID, requestID int64) error {
	return service.likeStorage.Create(requestID, profile1, profile2)
}

func (service *Service) GetMatches(profile types.UserID, requestID int64) ([]models.Card, error) {
	ids, err := service.likeStorage.GetMatch(requestID, profile)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	return service.GetProfile(ProfileGetParams{ID: ids}, requestID)
}
