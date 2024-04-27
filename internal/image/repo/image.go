package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"

	"main.go/config"
	"main.go/internal/image"
	. "main.go/internal/logs"
)

//

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

func GetImageRepo(config *config.DatabaseConfig) (*ImageStorage, error) {
	//logger := ctx.Value(Logg).(Log)

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	println(dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		println(err.Error())
	}
	if err = db.Ping(); err != nil {
		println(err.Error())
		//logger.Logger.Fatal(err)
	}

	postgreDb := ImageStorage{dbReader: db}

	go postgreDb.pingDb(config.Timer)
	return &postgreDb, nil
}

func (storage *ImageStorage) pingDb(timer uint32) {
	//logger := ctx.Value(Logg).(Log)
	for {
		err := storage.dbReader.Ping()
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Repo Profile db ping error ", err.Error())
		}

		time.Sleep(time.Duration(timer) * time.Second)
	}
}

func (storage *ImageStorage) Get(ctx context.Context, userID int64) ([]image.Image, error) {
	//logger := ctx.Value(Logg).(Log)
	//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Get request to person_image")
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
	//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Return ", len(images), " images")
	return images, nil
}

func (storage *ImageStorage) Add(ctx context.Context, image image.Image) error {
	logger := ctx.Value(Logg).(Log)
	query := "INSERT INTO person_image (person_id, image_url, cell_number) VALUES ($1, $2, $3)"

	_, err := storage.dbReader.Exec(query, image.UserId, image.Url, image.CellNumber)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}
	return nil
}

func (storage *ImageStorage) Delete(ctx context.Context, image image.Image) error {
	logger := ctx.Value(Logg).(Log)
	query := "DELETE FROM person_image WHERE person_id = $1 AND cell_number = $2"

	_, err := storage.dbReader.Exec(query, image.UserId, image.CellNumber)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Add img %w", err)
	}
	return nil
}
