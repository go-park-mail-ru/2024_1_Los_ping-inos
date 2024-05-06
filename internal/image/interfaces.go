package image

import (
	"context"
	"mime/multipart"
)

type (
	UseCase interface {
		GetImage(userID int64, cell string, ctx context.Context) (string, error)
		AddImage(userImage Image, img multipart.File, ctx context.Context) error
		DeleteImage(userImage Image, ctx context.Context) error
	}

	ImgStorage interface {
		Get(ctx context.Context, userID int64, cell string) ([]Image, error)
		Add(ctx context.Context, image Image) error
		Delete(ctx context.Context, image Image) error
	}
)
