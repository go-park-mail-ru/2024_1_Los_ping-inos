package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	auth "main.go/internal/auth/proto"
	"main.go/internal/image"
	"main.go/internal/image/repo"
)

const (
	vkCloudHotboxEndpoint = "https://hb.ru-msk.vkcs.cloud"
	defaultRegion         = "ru-msk"
)

type UseCase struct {
	imageStorage image.ImgStorage
	client       auth.AuthHandlClient
}

func NewImageUseCase(istore image.ImgStorage) *UseCase {
	return &UseCase{
		imageStorage: istore,
	}
}

func GetClient(port string) (auth.AuthHandlClient, error) {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc connect err: %w", err)
	}

	client := auth.NewAuthHandlClient(conn)

	return client, nil
}

func GetCore(cfg_sql string) (*UseCase, error) {
	grpclient, err := GetClient(":50051")
	images, err := repo.GetImageRepo(cfg_sql)

	if err != nil {
		return nil, err
	}

	core := UseCase{
		imageStorage: images,
		client:       grpclient,
	}
	return &core, nil
}

func (service *UseCase) GetImage(userID int64, cell string, ctx context.Context) (string, error) {
	images, err := service.imageStorage.Get(ctx, userID, cell)
	if err != nil {
		return "", err
	}

	if images == "" {
		return "", errors.New("no images for user with such sessionID")
	}

	return images, err
}

func (service *UseCase) AddImage(userImage image.Image, img multipart.File, ctx context.Context) error {

	err := service.imageStorage.Add(ctx, userImage, img)
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
