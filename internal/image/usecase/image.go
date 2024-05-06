package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	awsUpload "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	serviceUpload "github.com/aws/aws-sdk-go/service/s3"
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
	var url string

	for _, img := range images {
		objectKey := img.Url
		println("THIS IS OBJECT KEY", objectKey)
		req, err = presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(lifeTimeSeconds * int64(time.Second))
		})
		url = req.URL
		println(url)
	}

	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v. Error: %v\n", bucketName, err)
	}

	return url, err
}

func (service *UseCase) AddImage(userImage image.Image, img multipart.File, ctx context.Context) error {

	err := service.imageStorage.Add(ctx, userImage)
	if err != nil {
		return err
	}

	sess, err := session.NewSession(&awsUpload.Config{
		Region: aws.String("ru-msk"),
	})
	if err != nil {
		return err
	}

	svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	bucket := "los_ping"

	params := &serviceUpload.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(userImage.FileName),
		Body:   img,
		ACL:    aws.String("public-read"),
	}

	_, err = svc.PutObject(params)
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

	sess, err := session.NewSession(&awsUpload.Config{
		Region: aws.String("ru-msk"),
	})
	if err != nil {
		return err
	}

	svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	bucket := "los_ping"
	key := fmt.Sprint(userImage.UserId) + "/" + userImage.CellNumber + "/"

	input := &serviceUpload.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	}
	result, err := svc.ListObjectsV2(input)
	if err != nil {
		log.Fatalf("Unable to list objects in directory %q, %v\n", key, err)
	}

	for _, obj := range result.Contents {
		if _, err := svc.DeleteObject(&serviceUpload.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    obj.Key,
		}); err != nil {
			log.Fatalf("Unable to delete object %q from bucket %q, %v\n", key, bucket, err)
		} else {
			log.Printf("Object %q deleted from bucket %q\n", key, bucket)
		}
	}

	return nil
}
