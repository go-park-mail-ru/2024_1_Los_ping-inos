package delivery

import (
	"errors"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"math/rand"
	"net/http"
	"time"

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

	prometheus.MustRegister(
		image.TotalHits,
		image.HitDuration,
	)

	err := http.ListenAndServe(":8082", MetricTimeMiddleware(deliver.mx))
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

	grpcConn, err := grpc.Dial("auth:50051", grpc.WithInsecure())
	if err != nil {
		print("fock off")
	}
	authManager := gen.NewAuthHandlClient(grpcConn)

	api.mx.Handle("/metrics", promhttp.Handler())

	api.mx.Handle(apiPath+"getImage", requests.RequestIDMiddleware(
		requests.AllowedMethodMiddleware(
			requests.IsAuthenticatedMiddleware(http.HandlerFunc(api.GetImageHandler()), authManager), hashset.New("GET")),
		"get images", logger))

	api.mx.Handle(apiPath+"addImage", requests.RequestIDMiddleware(
		requests.AllowedMethodMiddleware(
			requests.IsAuthenticatedMiddleware(http.HandlerFunc(api.AddImageHandler()), authManager), hashset.New("POST")),
		"username (/me)", logger))

	api.mx.Handle(apiPath+"deleteImage", requests.RequestIDMiddleware(
		requests.AllowedMethodMiddleware(
			requests.IsAuthenticatedMiddleware(http.HandlerFunc(api.DeleteImageHandler()), authManager), hashset.New("POST")),
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

		cell := request.FormValue("cell")
		println(cell)

		userId := int64(request.Context().Value(RequestUserID).(types.UserID))

		images, err := deliver.useCase.GetImage(userId, cell, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		requests.SendSimpleResponse(respWriter, request, http.StatusOK, images)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent image")
	}
}

func (deliver *ImageHandler) AddImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)

		err := request.ParseMultipartForm(10 << 20)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		cell := request.FormValue("cell")
		img, handler, err := request.FormFile("image")
		if err != nil && errors.Is(err, http.ErrMissingFile) {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
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
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, "Wrong format")
			return
		}

		userId := int64(request.Context().Value(RequestUserID).(types.UserID))
		filename := fmt.Sprint(userId) + "/" + fmt.Sprint(cell) + "/" + fmt.Sprint(rand.Int()) + handler.Filename
		objectURL := "https://los_ping.hb.ru-msk.vkcs.cloud/" + filename
		userImage := image.Image{
			UserId:     userId,
			Url:        objectURL,
			CellNumber: cell,
			FileName:   filename,
		}

		err = deliver.useCase.AddImage(userImage, img, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("image added")
		requests.SendSimpleResponse(respWriter, request, http.StatusOK, objectURL)

	}
}

func (deliver *ImageHandler) DeleteImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		logger := request.Context().Value(Logg).(Log)
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("delete image")
		userId := int64(request.Context().Value(RequestUserID).(types.UserID))
		var r image.ImgRequest

		body, err := io.ReadAll(request.Body)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("bad body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		err = easyjson.Unmarshal(body, &r)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't unmarshal body: ", err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		userImage := image.Image{
			UserId:     userId,
			CellNumber: r.CellNumber,
		}

		err = deliver.useCase.DeleteImage(userImage, request.Context())
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("image added")
		requests.SendSimpleResponse(respWriter, request, http.StatusOK, "")
	}
}

func MetricTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next.ServeHTTP(respWriter, request)

		end := time.Since(start)
		path := request.URL.Path
		if path != "/metrics" {
			image.TotalHits.WithLabelValues().Inc()
			image.HitDuration.WithLabelValues(request.Method, path).Set(float64(end.Milliseconds()))
		}
	})
}
