package storage

import (
	"database/sql"

	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	models "main.go/db"
	. "main.go/internal/logs"
)

type InterestStorage struct {
	dbReader *sql.DB
}

func NewInterestStorage(dbReader *sql.DB) *InterestStorage {
	return &InterestStorage{
		dbReader: dbReader,
	}
}

func (storage *InterestStorage) Get(requestID int64) ([]*models.Interest, error) { // TODO добавить фильтры, когда продумаем интересы
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Select("*").
		From(InterestTableName).
		RunWith(storage.dbReader)

	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db update request to ", InterestTableName)
	rows, err := query.Query()
	defer rows.Close()

	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't query: ", err.Error())
		return nil, err
	}

	interests := make([]*models.Interest, 0)
	for rows.Next() {
		interest := &models.Interest{}
		err = rows.Scan(&interest.ID, &interest.Name)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't scan row: ", err.Error())
			return nil, err
		}

		interests = append(interests, interest)
	}

	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("return interests")
	return interests, nil
}
