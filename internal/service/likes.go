package service

import "main.go/internal/types"

func (service *Service) CreateLike(profile1, profile2 types.UserID, requestID int64) error {
	return service.likeStorage.Create(requestID, profile1, profile2)
}
