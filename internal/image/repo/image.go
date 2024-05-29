package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	_ "github.com/lib/pq"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	awsUpload "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	serviceUpload "github.com/aws/aws-sdk-go/service/s3"

	"main.go/internal/image"
	. "main.go/internal/logs"
)

const (
	vkCloudHotboxEndpoint = "https://hb.ru-msk.vkcs.cloud"
	defaultRegion         = "ru-msk"
)

const (
	personImageFields = "person_id, image_url, cell_number"
)

type ImageStorage struct {
	dbReader *sql.DB
}

func NewImageStorage(dbReader *sql.DB) *ImageStorage {
	return &ImageStorage{
		dbReader: dbReader,
	}
}

func GetImageRepo(config string) (*ImageStorage, error) {
	db, err := sql.Open("postgres", config)
	if err != nil {
		println(err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	if err = db.Ping(); err != nil {
		println(err.Error())
		//logger.Logger.Fatal(err)
	}

	postgreDb := ImageStorage{dbReader: db}

	go postgreDb.pingDb(50)
	return &postgreDb, nil
}

func (storage *ImageStorage) pingDb(timer uint32) {
	//logger := ctx.Value(Logg).(Log)
	for {
		err := storage.dbReader.Ping()
		if err != nil {
			println(err.Error())
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Repo Profile db ping error ", err.Error())
		}

		time.Sleep(time.Duration(timer) * time.Second)
	}
}

func (storage *ImageStorage) Get(ctx context.Context, userID int64, cell string) (string, error) {
	//logger := ctx.Value(Logg).(Log)
	//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Get request to person_image")
	var images []image.Image

	query := "SELECT " + personImageFields + " FROM person_image WHERE person_id = $1 AND cell_number = $2"

	stmt, err := storage.dbReader.Prepare(query) // using prepared statement
	if err != nil {
		return "", err
	}
	rows, err := stmt.Query(userID, cell)
	//rows, err := storage.dbReader.Query(query, userID, cell)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var img image.Image

		err := rows.Scan(&img.UserId, &img.Url, &img.CellNumber)
		if err != nil {
			return "", err
		}

		images = append(images, img)
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

	//var req *v4.PresignedHTTPRequest
	//var url string
	var obj string

	for _, img := range images {
		objectKey := img.Url
		println("THIS IS OBJECT KEY", objectKey)
		_, err = presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(lifeTimeSeconds * int64(time.Second))
		})
		if err != nil {
			println(err.Error())
			return "", err
		}
		//url = req.URL
		//println(url)
		obj = img.Url
		return obj, nil
	}
	//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Return ", len(images), " images")
	return obj, nil
}

func (storage *ImageStorage) Add(ctx context.Context, image image.Image, img multipart.File) error {
	logger := ctx.Value(Logg).(Log)
	query := "INSERT INTO person_image (person_id, image_url, cell_number) VALUES ($1, $2, $3) ON CONFLICT (person_id, cell_number) DO UPDATE SET image_url = EXCLUDED.image_url;"

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("hehe ", image.UserId, image.CellNumber, image.Url)
	stmt, err := storage.dbReader.Prepare(query) // using prepared statement
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}

	_, err = stmt.Exec(image.UserId, image.Url, image.CellNumber)
	//_, err := storage.dbReader.Exec(query, image.UserId, image.Url, image.CellNumber)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
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
		Key:    aws.String(image.FileName),
		Body:   img,
		ACL:    aws.String("public-read"),
	}

	_, err = svc.PutObject(params)
	if err != nil {
		return err
	}
	return nil
}

func (storage *ImageStorage) Delete(ctx context.Context, image image.Image) error {
	logger := ctx.Value(Logg).(Log)
	query := "DELETE FROM person_image WHERE person_id = $1 AND cell_number = $2"

	stmt, err := storage.dbReader.Prepare(query) // using prepared statement
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}
	_, err = stmt.Exec(image.UserId, image.CellNumber)
	//_, err := storage.dbReader.Exec(query, image.UserId, image.CellNumber)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}

	sess, err := session.NewSession(&awsUpload.Config{
		Region: aws.String("ru-msk"),
	})
	if err != nil {
		return err
	}

	svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
	bucket := "los_ping"
	key := fmt.Sprint(image.UserId) + "/" + image.CellNumber + "/"

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
