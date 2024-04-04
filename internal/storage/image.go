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

func (storage *ImageStorage) Get(requestID int64, person models.Person) (*models.Image, error) {
	//stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	// query := stBuilder.
	// 	Select("photo").
	// 	From(PersonTableName).
	// 	Where(qb.Eq{"id": person.ID}).
	// 	RunWith(storage.dbReader)

	// Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db get request to ", PersonTableName)
	// rows, err := query.Query()

	// if err != nil {
	// 	Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't query: ", err.Error())
	// 	return nil, err
	// }
	// defer rows.Close()

	// image := &models.Image{}
	// err = rows.Scan(&image.Url)

	// if err != nil {
	// 	Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't scan person: ", err.Error())
	// 	return nil, err
	// }

	// Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db returning records")
	// return image, nil

	//image := qb.Select("photo").From(PersonTableName).Where(qb.Eq{"id": person.ID})

	imageItem := &models.Image{}

	//print(person.SessionID)

	err := storage.dbReader.QueryRow(
		`SELECT id, photo FROM person
				WHERE person.session_id = $1`, person.SessionID).Scan(&imageItem.UserId, &imageItem.Url)

	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't scan person: ", err.Error())
		return nil, err
	}

	return imageItem, nil

}

func (storage *ImageStorage) Add(requestID int64, image models.Image) error {
	_, err := storage.dbReader.Exec(
		"UPDATE person SET photo = $1 WHERE session_id = $2", image.Url, image.UserId)

	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}
	//Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db added img ", PersonTableName)
	return nil
}
