package storage

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	//. "main.go/config"
	models "main.go/db"
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

func (storage *ImageStorage) Get(requestID int64, userID int64) ([]models.Image, error) {
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("Get request to person_image")
	var images []models.Image

	query := "SELECT" + personImageFields + " FROM person_image WHERE person_id = $1"

	rows, err := storage.dbReader.Query(query, userID)
	if err != nil {
		return []models.Image{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var image models.Image

		err := rows.Scan(&image.UserId, &image.Url, &image.CellNumber)
		if err != nil {
			return []models.Image{}, err
		}

		images = append(images, image)
	}
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("Return ", len(images), " images")
	return images, nil
}

func (storage *ImageStorage) Add(requestID int64, image models.Image) error {
	query := "INSERT INTO person_image (person_id, image_url, cell_number) VALUES ($1, $2, $3)"

	_, err := storage.dbReader.Exec(query, image.UserId, image.Url, image.CellNumber)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}
	return nil
}

func (storage *ImageStorage) Delete(requestID int64, image models.Image) error {
	query := "DELETE FROM person_image WHERE person_id = $1 AND cell_number = $2"

	_, err := storage.dbReader.Exec(query, image.UserId, image.CellNumber)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}
	return nil
}
