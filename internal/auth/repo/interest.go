package repo

import (
	"context"
	"database/sql"

	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"main.go/internal/auth"
	. "main.go/internal/logs"
	"main.go/internal/types"
)

const (
	PersonInterestTableName = "person_interest"
	interestFields          = "id, name"
	InterestTableName       = "interest"
	personInterestFields    = "person_id, interest_id"
)

type InterestStorage struct {
	dbReader *sql.DB
}

func NewInterestStorage(dbReader *sql.DB) *InterestStorage {
	return &InterestStorage{
		dbReader: dbReader,
	}
}

func (storage *InterestStorage) CreatePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error {
	//logger := ctx.Value(Logg).(Log)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db add request to ", PersonInterestTableName)

	for i := range interestID {
		query := stBuilder.
			Insert(PersonInterestTableName).
			Columns("person_id", "interest_id").
			Values(personID, interestID[i]).
			RunWith(storage.dbReader)

		rows, err := query.Query()
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db insert can't query: ", err.Error())
			return err
		}
		rows.Close()
	}

	return nil
}

func (storage *InterestStorage) Get(ctx context.Context, filter *auth.InterestGetFilter) ([]*auth.Interest, error) {
	logger := ctx.Value(Logg).(Log)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	whereMap := qb.And{}

	if filter != nil && filter.ID == nil && filter.Name == nil {
		return nil, nil
	}

	if filter == nil {
		filter = &auth.InterestGetFilter{}
	}

	processInterestIDFilter(filter, &whereMap)
	processInterestNameFilter(filter, &whereMap)

	query := stBuilder.
		Select(interestFields).
		From(InterestTableName).
		Where(whereMap).
		RunWith(storage.dbReader)

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to ", InterestTableName)
	rows, err := query.Query()
	defer rows.Close()

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return nil, err
	}

	interests := make([]*auth.Interest, 0)
	for rows.Next() {
		interest := &auth.Interest{}
		err = rows.Scan(&interest.ID, &interest.Name)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't scan row: ", err.Error())
			return nil, err
		}

		interests = append(interests, interest)
	}

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("return ", len(interests), " interests")
	return interests, nil
}

func processInterestIDFilter(filter *auth.InterestGetFilter, whereMap *qb.And) {
	if filter.ID != nil {
		*whereMap = append(*whereMap, qb.Eq{"id": filter.ID})
	}
}

func processInterestNameFilter(filter *auth.InterestGetFilter, whereMap *qb.And) {
	if filter.Name != nil {
		*whereMap = append(*whereMap, qb.Eq{"name": filter.Name})
	}
}

func (storage *InterestStorage) GetPersonInterests(ctx context.Context, personID types.UserID) ([]*auth.Interest, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to ", PersonInterestTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	query := stBuilder.
		Select(personInterestFields).
		From(PersonInterestTableName).
		Where(qb.Eq{"person_id": personID}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't query: ", err.Error())
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
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db person_interest can't scan: ", err.Error())
			return nil, err
		}
		ids = append(ids, interestID)
	}
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db got ", len(ids), " interest ids")
	return storage.Get(ctx, &auth.InterestGetFilter{ID: ids})
}

func (storage *InterestStorage) DeletePersonInterests(ctx context.Context, personID types.UserID, interestID []types.InterestID) error {
	logger := ctx.Value(Logg).(Log)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db delete request to ", PersonInterestTableName)
	query := stBuilder.
		Delete(PersonInterestTableName).
		Where(qb.And{qb.Eq{"person_id": personID}, qb.Eq{"interest_id": interestID}}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	defer rows.Close()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db delete can't query: ", err.Error())
		return err
	}
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db person interest deleted")
	return nil
}
