package service

import "main.go/internal/types"

func (service *Service) CreateLike(profile1, profile2 types.UserID, requestID int64) error {
	return service.likeStorage.Create(requestID, profile1, profile2)
}

func (service *Service) GetMatches(profile types.UserID, requestID int64) (string, error) {
	ids, err := service.likeStorage.GetMatch(requestID, profile)
	if err != nil {
		return "", err
	}

	return service.GetProfile(ProfileGetParams{ID: ids}, requestID)
}
