package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	models "main.go/db"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsUpload "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	serviceUpload "github.com/aws/aws-sdk-go/service/s3"
	. "main.go/config"
)

const (
	vkCloudHotboxEndpoint = "https://hb.ru-msk.vkcs.cloud"
	defaultRegion         = "ru-msk"
)

func (deliver *Deliver) GetImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		requestID := request.Context().Value(RequestID).(int64)

		userId := int64(request.Context().Value(RequestUserID).(types.UserID))

		images, err := deliver.serv.GetImage(userId, requestID)
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
		lifeTimeSeconds := int64(60)

		var req *v4.PresignedHTTPRequest
		//urls := make([]string, 0)
		url := make(map[string][]map[string]string)

		for _, image := range images {
			objectKey := image.Url
			req, err = presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
			}, func(opts *s3.PresignOptions) {
				opts.Expires = time.Duration(lifeTimeSeconds * int64(time.Second))
			})
			newImage := map[string]string{
				"cell": image.CellNumber,
				"url":  req.URL,
			}
			//urls = append(urls, req.URL)
			url["photo"] = append(url["photo"], newImage)
		}

		if err != nil {
			log.Printf("Couldn't get a presigned request to get %v. Error: %v\n", bucketName, err)
		}

		requests.SendResponse(respWriter, request, http.StatusOK, url)
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("sent image")
	}
}

func (deliver *Deliver) AddImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		requestID := request.Context().Value(RequestID).(int64)
		//var r requests.ImageAddRequest

		err := request.ParseMultipartForm(10 << 20)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		cell := request.FormValue("cell")

		image, handler, err := request.FormFile("image")
		if err != nil && errors.Is(err, http.ErrMissingFile) {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}
		defer image.Close()

		fileType := handler.Header.Get("Content-Type")

		isValidImage := false
		for _, validType := range types.ValidImageTypes {
			if fileType == validType {
				isValidImage = true
				break
			}
		}

		if !isValidImage {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("wrong format")
			requests.SendResponse(respWriter, request, http.StatusBadRequest, "Wrong format")
			return
		}

		userId := int64(request.Context().Value(RequestUserID).(types.UserID))
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		filename := fmt.Sprint(userId) + "/" + fmt.Sprint(cell) + "/" + fmt.Sprint(rand.Int()) + handler.Filename
		objectURL := "https://los_ping.hb.ru-msk.vkcs.cloud/" + filename

		sess, err := session.NewSession(&awsUpload.Config{
			Region: aws.String("ru-msk"),
		})
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			return
		}

		userImage := models.Image{
			UserId:     userId,
			Url:        objectURL,
			CellNumber: cell,
		}

		err = deliver.serv.AddImage(userImage, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
		bucket := "los_ping"

		params := &serviceUpload.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
			Body:   image,
			ACL:    aws.String("public-read"),
		}

		_, err = svc.PutObject(params)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			return
		}

		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("image added")
		requests.SendResponse(respWriter, request, http.StatusOK, objectURL)

	}
}

func (deliver *Deliver) DeleteImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		requestID := deliver.nextRequest()
		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("delete image")
		userId := int64(request.Context().Value(RequestUserID).(types.UserID))
		var r requests.ImageRequest

		body, err := io.ReadAll(request.Body)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("bad body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		userImage := models.Image{
			UserId:     userId,
			CellNumber: r.CellNumber,
		}

		err = deliver.serv.DeleteImage(userImage, requestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		sess, err := session.NewSession(&awsUpload.Config{
			Region: aws.String("ru-msk"),
		})
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
			return
		}

		svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
		bucket := "los_ping"
		key := fmt.Sprint(userId) + "/" + r.CellNumber + "/"

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

		Log.WithFields(logrus.Fields{RequestID: requestID}).Info("image added")
		requests.SendResponse(respWriter, request, http.StatusOK, nil)
	}
}
