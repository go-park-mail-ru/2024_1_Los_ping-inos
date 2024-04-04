package storage

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	//. "main.go/config"
	models "main.go/db"
	. "main.go/internal/logs"
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
	var images []models.Image

	query := "SELECT * FROM person_image WHERE person_id = $1"

	rows, err := storage.dbReader.Query(query, userID)
	if err != nil {
		return []models.Image{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var image models.Image

		err := rows.Scan(&image.UserId, &image.Url)
		if err != nil {
			return []models.Image{}, err
		}

		images = append(images, image)
	}
	return images, nil
}

func (storage *ImageStorage) Add(requestID int64, image models.Image) error {
	_, err := storage.dbReader.Exec(
		"UPDATE person SET photo = $1 WHERE session_id = $2", image.Url, image.UserId)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}

	_, err = storage.dbReader.Exec(
		"INSERT INTO image (url) VALUES ($1)", image.Url)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}

	query := "INSERT INTO person_image (person_id, image_url) VALUES ($1, $2)"

	_, err = storage.dbReader.Exec(query, image.UserId, image.Url)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}
	return nil
}
