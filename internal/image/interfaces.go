package image

import (
	"context"
)

type (
	UseCase interface {
		GetImage(userID int64, ctx context.Context) ([]Image, error)
		AddImage(userImage Image, ctx context.Context) error
		DeleteImage(userImage Image, ctx context.Context) error
	}

	ImgStorage interface {
		Get(ctx context.Context, userID int64) ([]Image, error)
		Add(ctx context.Context, image Image) error
		Delete(ctx context.Context, image Image) error
	}
)
