package repo

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"main.go/internal/auth"
	. "main.go/internal/logs"
	"strconv"
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

func (storage *ImageStorage) Get(ctx context.Context, userID int64) ([]auth.Image, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Get request to person_image")

	query := "SELECT " + personImageFields + " FROM person_image WHERE person_id = $1"

	rows, err := storage.dbReader.Query(query, userID)
	if err != nil {
		return []auth.Image{}, err
	}
	defer rows.Close()

	images := make([]auth.Image, 5)
	for i := 0; i < 5; i++ {
		images[i] = auth.Image{CellNumber: strconv.Itoa(i), Url: ""}
	}

	for rows.Next() {
		image := auth.Image{}
		err = rows.Scan(&image.UserId, &image.Url, &image.CellNumber)
		if err != nil {
			return []auth.Image{}, err
		}
		cell, _ := strconv.Atoi(image.CellNumber)
		images[cell] = auth.Image{image.UserId, image.Url, image.CellNumber}
	}
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Return ", len(images), " images")
	return images, nil
}
