package usecase

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"main.go/config"
	auth "main.go/internal/auth/proto"
	"main.go/internal/image"
	"main.go/internal/image/repo"
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

func GetCore(cfg_sql *config.DatabaseConfig) (*UseCase, error) {
	grpclient, err := GetClient(cfg_sql.GrpcPort)
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
