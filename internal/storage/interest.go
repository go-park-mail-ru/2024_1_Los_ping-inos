package storage

import (
	"database/sql"

	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	models "main.go/db"
	. "main.go/internal/logs"
	"main.go/internal/types"
)

type InterestStorage struct {
	dbReader *sql.DB
}

func NewInterestStorage(dbReader *sql.DB) *InterestStorage {
	return &InterestStorage{
		dbReader: dbReader,
	}
}

func processInterestIDFilter(filter *models.InterestGetFilter, whereMap *qb.And) {
	if filter.ID != nil {
		*whereMap = append(*whereMap, qb.Eq{"id": filter.ID})
	}
}

func processInterestNameFilter(filter *models.InterestGetFilter, whereMap *qb.And) {
	if filter.Name != nil {
		*whereMap = append(*whereMap, qb.Eq{"name": filter.Name})
	}
}

func (storage *InterestStorage) Get(requestID int64, filter *models.InterestGetFilter) ([]*models.Interest, error) {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	whereMap := qb.And{}

	if filter != nil && filter.ID == nil && filter.Name == nil {
		return nil, nil
	}

	if filter == nil {
		filter = &models.InterestGetFilter{}
	}

	processInterestIDFilter(filter, &whereMap)
	processInterestNameFilter(filter, &whereMap)

	query := stBuilder.
		Select("*").
		From(InterestTableName).
		Where(whereMap).
		RunWith(storage.dbReader)

	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db get request to ", InterestTableName)
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

	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("return ", len(interests), " interests")
	return interests, nil
}

func (storage *InterestStorage) GetPersonInterests(requestID int64, personID types.UserID) ([]*models.Interest, error) {
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db get request to ", PersonInterestTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	query := stBuilder.Select("*").
		From(PersonInterestTableName).
		Where(qb.Eq{"person_id": personID}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	var (
		ids        []types.InterestID
		personsID  types.UserID
		interestID types.InterestID
	)
	for rows.Next() {
		err = rows.Scan(&personsID, &interestID)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db person_interest can't scan: ", err.Error())
			return nil, err
		}
		ids = append(ids, interestID)
	}
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db got ", len(ids), " interest ids")
	return storage.Get(requestID, &models.InterestGetFilter{ID: ids})
}

func (storage *InterestStorage) CreatePersonInterests(requestID int64, personID types.UserID, interestID []types.InterestID) error {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db add request to ", PersonInterestTableName)

	for i := range interestID {
		query := stBuilder.
			Insert(PersonInterestTableName).
			Columns("person_id", "interest_id").
			Values(personID, interestID[i]).
			RunWith(storage.dbReader)

		rows, err := query.Query()
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db insert can't query: ", err.Error())
			rows.Close()
			return err
		}
		rows.Close()
	}

	return nil
}

func (storage *InterestStorage) DeletePersonInterests(requestID int64, personID types.UserID, interestID []types.InterestID) error {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db delete request to ", PersonInterestTableName)
	query := stBuilder.
		Delete(PersonInterestTableName).
		Where(qb.And{qb.Eq{"person_id": personID}, qb.Eq{"interest_id": interestID}}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	defer rows.Close()
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db delete can't query: ", err.Error())
		return err
	}
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db person interest deleted")
	return nil
}
