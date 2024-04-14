package repo

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"main.go/internal/feed"
	. "main.go/internal/logs"
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

func (storage *ImageStorage) Get(ctx context.Context, userID int64) ([]feed.Image, error) {
	logger := ctx.Value(Logg).(*Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Get request to person_image")
	var images []feed.Image

	query := "SELECT " + personImageFields + " FROM person_image WHERE person_id = $1"

	rows, err := storage.dbReader.Query(query, userID)
	if err != nil {
		return []feed.Image{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var image feed.Image

		err := rows.Scan(&image.UserId, &image.Url, &image.CellNumber)
		if err != nil {
			return []feed.Image{}, err
		}

		images = append(images, image)
	}
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Return ", len(images), " images")
	return images, nil
}
