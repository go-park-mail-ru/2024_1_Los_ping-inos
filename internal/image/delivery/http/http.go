package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsUpload "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	serviceUpload "github.com/aws/aws-sdk-go/service/s3"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	. "main.go/config"
	gen "main.go/internal/auth/proto"
	"main.go/internal/image"
	"main.go/internal/image/usecase"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"
)

type ImageHandler struct {
	useCase image.UseCase
	mx      *http.ServeMux
}

func (deliver *ImageHandler) ListenAndServe() error {
	// server := http.Server{
	// 	Addr:         cfg.Host + cfg.Port,
	// 	Handler:      deliver.mx,
	// 	ReadTimeout:  cfg.Timeout * time.Second,
	// 	WriteTimeout: cfg.Timeout * time.Second,
	// }

	//logger.Logger.Infof("started auth http server at %v", server.Addr)
	//	fmt.Printf("started auth http server at %v\n", server.Addr)
	err := http.ListenAndServe(":8090", deliver.mx)
	if err != nil {
		//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
		return fmt.Errorf("listen and serve error: %w", err)
	}

	return nil
}

func GetApi(c *usecase.UseCase, logger Log) *ImageHandler {
	api := &ImageHandler{
		useCase: c,
		mx:      http.NewServeMux(),
	}
	var apiPath = "/api/v1/"

	println("This is api path", apiPath)

	grpcConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		print("fock off")
	}
	authManager := gen.NewAuthHandlClient(grpcConn)

	api.mx.Handle(apiPath+"getImage", requests.RequestIDMiddleware(
		requests.AllowedMethodMiddleware(
			requests.IsAuthenticatedMiddleware(http.HandlerFunc(api.GetImageHandler()), authManager), hashset.New("GET")),
		"get images", logger))

	api.mx.Handle(apiPath+"addImage", requests.RequestIDMiddleware(
		requests.AllowedMethodMiddleware(
			http.HandlerFunc(api.AddImageHandler()), hashset.New("POST")),
		"username (/me)", logger))

	api.mx.Handle(apiPath+"deleteImage", requests.RequestIDMiddleware(
		requests.AllowedMethodMiddleware(
			http.HandlerFunc(api.DeleteImageHandler()), hashset.New("POST")),
		"delete image", logger))

	return api
}

func NewImageDelivery(uc image.UseCase) *ImageHandler {
	return &ImageHandler{
		useCase: uc,
	}
}

const (
	vkCloudHotboxEndpoint = "https://hb.ru-msk.vkcs.cloud"
	defaultRegion         = "ru-msk"
)

func (deliver *ImageHandler) GetImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)

		//var requestBody requests.ImageRequest

		cell := request.FormValue("cell")
		println(cell)

		userId := int64(request.Context().Value(RequestUserID).(types.UserID))

		images, err := deliver.useCase.GetImage(userId, cell, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		requests.SendResponse(respWriter, request, http.StatusOK, images)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent image")
	}
}

func (deliver *ImageHandler) AddImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)

		err := request.ParseMultipartForm(10 << 20)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		cell := request.FormValue("cell")

		img, handler, err := request.FormFile("image")
		if err != nil && errors.Is(err, http.ErrMissingFile) {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}
		defer img.Close()

		fileType := handler.Header.Get("Content-Type")

		isValidImage := false
		for _, validType := range types.ValidImageTypes {
			if fileType == validType {
				isValidImage = true
				break
			}
		}

		if !isValidImage {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("wrong format")
			requests.SendResponse(respWriter, request, http.StatusBadRequest, "Wrong format")
			return
		}

		//userId := int64(request.Context().Value(RequestUserID).(types.UserID))
		userId := int64(2)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		filename := fmt.Sprint(userId) + "/" + fmt.Sprint(cell) + "/" + fmt.Sprint(rand.Int()) + handler.Filename
		objectURL := "https://los_ping.hb.ru-msk.vkcs.cloud/" + filename

		sess, err := session.NewSession(&awsUpload.Config{
			Region: aws.String("ru-msk"),
		})
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			return
		}

		userImage := image.Image{
			UserId:     userId,
			Url:        objectURL,
			CellNumber: cell,
		}

		err = deliver.useCase.AddImage(userImage, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		svc := serviceUpload.New(sess, awsUpload.NewConfig().WithEndpoint(vkCloudHotboxEndpoint).WithRegion(defaultRegion))
		bucket := "los_ping"

		params := &serviceUpload.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
			Body:   img,
			ACL:    aws.String("public-read"),
		}

		_, err = svc.PutObject(params)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			return
		}

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("image added")
		requests.SendResponse(respWriter, request, http.StatusOK, objectURL)

	}
}

func (deliver *ImageHandler) DeleteImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("delete image")
		//userId := int64(request.Context().Value(RequestUserID).(types.UserID))
		userId := int64(2)
		var r requests.ImageRequest

		body, err := io.ReadAll(request.Body)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("bad body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		userImage := image.Image{
			UserId:     userId,
			CellNumber: r.CellNumber,
		}

		err = deliver.useCase.DeleteImage(userImage, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		sess, err := session.NewSession(&awsUpload.Config{
			Region: aws.String("ru-msk"),
		})
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
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

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("image added")
		requests.SendResponse(respWriter, request, http.StatusOK, nil)
	}
}
