package storage

import (
	"context"
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	. "main.go/internal/logs"
	"main.go/internal/profile"
	"main.go/internal/types"
)

const (
	likeFields = "person_one_id, person_two_id"
)

type LikeStorage struct {
	dbReader *sql.DB
}

func NewLikeStorage(dbReader *sql.DB) *LikeStorage {
	return &LikeStorage{
		dbReader: dbReader,
	}
}

func (storage *LikeStorage) Get(ctx context.Context, filter *profile.LikeGetFilter) ([]types.UserID, error) {
	logger := ctx.Value(Logg).(*Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to ", LikeTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	whereMap := qb.Or{}

	if filter == nil {
		filter = &profile.LikeGetFilter{}
	}

	if filter.Person1 != nil {
		whereMap = append(whereMap, qb.Eq{"person_one_id": filter.Person1})
	}

	query := stBuilder.
		Select(likeFields).
		From(LikeTableName).
		Where(whereMap).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	res := make([]types.UserID, 0)
	var tmp profile.Like
	for rows.Next() {
		err = rows.Scan(&tmp.Person1, &tmp.Person2)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db can't scan: ", err.Error())
			return nil, err
		}
		res = append(res, tmp.Person2)
	}

	return res, nil
}

func (storage *LikeStorage) Create(ctx context.Context, person1ID, person2ID types.UserID) error {
	logger := ctx.Value(Logg).(*Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db create request to ", LikeTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Insert(LikeTableName).
		Columns("person_one_id", "person_two_id").
		Values(person1ID, person2ID).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't create like: ", err.Error())
		return err
	}
	defer rows.Close()
	return nil
}

func (storage *LikeStorage) GetMatch(ctx context.Context, person1ID types.UserID) ([]types.UserID, error) {
	logger := ctx.Value(Logg).(*Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to ", LikeTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.Select("t1.person_two_id").
		From(LikeTableName + "t1").
		InnerJoin(LikeTableName + "t2 ON t1.person_one_id = t2.person_two_id AND t1.person_two_id = t2.person_one_id").
		Where(qb.Eq{"t1.person_one_id": person1ID}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	defer rows.Close()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't query:  ", err.Error())
		return nil, err
	}

	res := make([]types.UserID, 0)
	var scan types.UserID
	for rows.Next() {
		err = rows.Scan(&scan)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't scan:  ", err.Error())
			return nil, err
		}
		res = append(res, scan)
	}
	return res, nil
}
