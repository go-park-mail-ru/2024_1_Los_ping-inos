package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
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

	if images == nil {
		return "", errors.New("no images for user with such sessionID")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		print("Error loading default config: %v", err)
		os.Exit(0)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(vkCloudHotboxEndpoint)
		o.Region = defaultRegion
	})

	presigner := s3.NewPresignClient(client)
	bucketName := "los_ping"
	lifeTimeSeconds := int64(60)

	var req *v4.PresignedHTTPRequest
	//urls := make([]string, 0)
	var url string

	for _, img := range images {
		objectKey := img.Url
		println(objectKey)
		req, err = presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(lifeTimeSeconds * int64(time.Second))
		})
		//urls = append(urls, img.CellNumber, req.URL)
		url = req.URL
	}

	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v. Error: %v\n", bucketName, err)
	}

	return url, err
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
