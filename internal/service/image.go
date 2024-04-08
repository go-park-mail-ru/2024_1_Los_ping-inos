package service

import (
	"errors"

	models "main.go/db"
)

func (service *Service) GetImage(userID int64, requestID int64) ([]models.Image, error) {
	images, err := service.imageStorage.Get(requestID, userID)
	if err != nil {
		return nil, err
	}

	if images == nil {
		return nil, errors.New("no images for user with such sessionID")
	}

	return images, err
}

func (service *Service) AddImage(userImage models.Image, requestID int64) error {
	err := service.imageStorage.Add(requestID, userImage)
	if err != nil {
		return err
	}

	return nil
}
