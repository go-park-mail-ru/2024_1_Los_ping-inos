package delivery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	models "main.go/db"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsUpload "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	serviceUpload "github.com/aws/aws-sdk-go/service/s3"
)

const (
	vkCloudHotboxEndpoint = "https://hb.ru-msk.vkcs.cloud"
	defaultRegion         = "us-east-1"
)

func (deliver *Deliver) GetImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		requestID := deliver.nextRequest()
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("get user image")

		userSession, _ := request.Cookie("session_id")
		print(userSession.Value)
		print(" - THIS IS SESSION ID")

		image, err := deliver.serv.GetImage(userSession.Value, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
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
		objectKey := image
		lifeTimeSeconds := int64(60)

		fmt.Println("OBJECT KEY-", objectKey, "\n")

		req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(lifeTimeSeconds * int64(time.Second))
		})

		if err != nil {
			log.Printf("Couldn't get a presigned request to get %v:%v. Error: %v\n", bucketName, objectKey, err)
		}

		fmt.Printf("%s", req.URL)

		requests.SendResponse(respWriter, request, http.StatusOK, image)
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("sent image")
	}
}

func (deliver *Deliver) AddImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		requestID := deliver.nextRequest()
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("get upload image")

		err := request.ParseMultipartForm(10 << 20)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		image, handler, err := request.FormFile("image")
		if err != nil && errors.Is(err, http.ErrMissingFile) {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}
		defer image.Close()

		// fileType := handler.Header.Get("Content-Type")
		// isValidImage := false
		// for _, validType := range types.ValidImageTypes {
		// 	if fileType == validType {
		// 		isValidImage = true
		// 		break
		// 	}
		// }

		// print("777")

		// if err != nil || handler == nil || image == nil || !isValidImage {
		// 	Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
		// 	requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
		// 	return
		// }

		user_session, _ := request.Cookie("session_id")

		userId, err := deliver.serv.GetId(user_session.Value, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("this is user id -", userId)

		filename := fmt.Sprint(userId) + "/" + handler.Filename
		if err != nil && handler != nil && image != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		sess, err := session.NewSession(&awsUpload.Config{
			Region:      aws.String("ru-msk"),
			Credentials: credentials.NewStaticCredentials("jFMjTLNLjWRqR6uyTdZYkT", "6BHRrZdvVntY2hVgkdppMphZbLSj5YXyoVq4GTCzBuZk", ""),
		})
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			return
		}

		userImage := models.Image{
			UserId: user_session.Value,
			Url:    filename,
		}

		svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
		bucket := "los_ping"

		params := &serviceUpload.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
			Body:   image,
		}

		_, err = svc.PutObject(params)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			return
		}

		//fmt.Printf("data: %v\n", userImage)

		err = deliver.serv.AddImage(userImage, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("image added")
		requests.SendResponse(respWriter, request, http.StatusOK, nil)

	}
}
