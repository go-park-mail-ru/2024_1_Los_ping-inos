package usecase

import (
	"context"
	"errors"
	"main.go/internal/image"
)

type UseCase struct {
	imageStorage image.ImgStorage
}

func NewImageUseCase(istore image.ImgStorage) *UseCase {
	return &UseCase{
		imageStorage: istore,
	}
}

func (service *UseCase) GetImage(userID int64, ctx context.Context) ([]image.Image, error) {
	images, err := service.imageStorage.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if images == nil {
		return nil, errors.New("no images for user with such sessionID")
	}

	return images, err
}

func (service *UseCase) AddImage(userImage image.Image, ctx context.Context) error {
	err := service.imageStorage.Add(ctx, userImage)
	if err != nil {
		return err
	}

	return nil
}

func (service *UseCase) DeleteImage(userImage image.Image, ctx context.Context) error {
	err := service.imageStorage.Delete(ctx, userImage)
	if err != nil {
		return err
	}

	return nil
}
