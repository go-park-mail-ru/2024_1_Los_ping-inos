package storage

import (
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	models "main.go/db"
)

type InterestStorage struct {
	dbReader *sql.DB
}

func NewInterestStorage(dbReader *sql.DB) *InterestStorage {
	return &InterestStorage{
		dbReader: dbReader,
	}
}

func (storage *InterestStorage) Get() ([]*models.Interest, error) { // TODO добавить фильтры, когда продумаем интересы
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Select("*").
		From(InterestTableName).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	defer rows.Close()

	if err != nil {
		logrus.Info("can't read Interests: ", err.Error())
		return nil, err
	}

	interests := make([]*models.Interest, 0)
	for rows.Next() {
		interest := &models.Interest{}
		err = rows.Scan(&interest.ID, &interest.Name)
		if err != nil {
			logrus.Info("can't scan row ", err.Error())
			return nil, err
		}

		interests = append(interests, interest)
	}

	return interests, nil
}
