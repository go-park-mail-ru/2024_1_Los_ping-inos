package service

import (
	"context"
	"errors"

	models "main.go/db"
)

func (service *Service) GetImage(userID int64, ctx context.Context) ([]models.Image, error) {
	images, err := service.imageStorage.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if images == nil {
		return nil, errors.New("no images for user with such sessionID")
	}

	return images, err
}

func (service *Service) AddImage(userImage models.Image, ctx context.Context) error {
	err := service.imageStorage.Add(ctx, userImage)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) DeleteImage(userImage models.Image, ctx context.Context) error {
	err := service.imageStorage.Delete(ctx, userImage)
	if err != nil {
		return err
	}

	return nil
}
