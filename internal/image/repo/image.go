package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"main.go/internal/image"
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

func (storage *ImageStorage) Get(ctx context.Context, userID int64) ([]image.Image, error) {
	logger := ctx.Value(Logg).(*Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Get request to person_image")
	var images []image.Image

	query := "SELECT " + personImageFields + " FROM person_image WHERE person_id = $1"

	rows, err := storage.dbReader.Query(query, userID)
	if err != nil {
		return []image.Image{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var img image.Image

		err := rows.Scan(&img.UserId, &img.Url, &img.CellNumber)
		if err != nil {
			return []image.Image{}, err
		}

		images = append(images, img)
	}
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Return ", len(images), " images")
	return images, nil
}

func (storage *ImageStorage) Add(ctx context.Context, image image.Image) error {
	logger := ctx.Value(Logg).(*Log)
	query := "INSERT INTO person_image (person_id, image_url, cell_number) VALUES ($1, $2, $3)"

	_, err := storage.dbReader.Exec(query, image.UserId, image.Url, image.CellNumber)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}
	return nil
}

func (storage *ImageStorage) Delete(ctx context.Context, image image.Image) error {
	logger := ctx.Value(Logg).(*Log)
	query := "DELETE FROM person_image WHERE person_id = $1 AND cell_number = $2"

	_, err := storage.dbReader.Exec(query, image.UserId, image.CellNumber)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}
	return nil
}
