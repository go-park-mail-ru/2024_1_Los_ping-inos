package service

import (
	"errors"

	models "main.go/db"
)

func (service *Service) GetImage(sessionID string, requestID int64) (string, error) {
	//print(sessionID, "SEESION IS")
	image, err := service.imageStorage.Get(requestID, models.Person{SessionID: sessionID})
	if err != nil {
		return "", err
	}

	if image == nil {
		return "", errors.New("no images for user with such sessionID")
	}

	return image.Url, err
}

func (service *Service) AddImage(userImage models.Image, requestID int64) error {
	err := service.imageStorage.Add(requestID, userImage)
	if err != nil {
		return err
	}

	return nil
}
