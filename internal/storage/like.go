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

type LikeStorage struct {
	dbReader *sql.DB
}

func NewLikeStorage(dbReader *sql.DB) *LikeStorage {
	return &LikeStorage{
		dbReader: dbReader,
	}
}

func (storage *LikeStorage) Get(requestID int64, filter *models.LikeGetFilter) ([]*models.Like, error) {
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db get request to ", LikeTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	whereMap := qb.And{}

	if filter == nil {
		filter = &models.LikeGetFilter{}
	}

	if filter.Person1 != nil {
		whereMap = append(whereMap, qb.Eq{"person_one_id": filter.Person1})
	}
	if filter.Person2 != nil {
		whereMap = append(whereMap, qb.Eq{"person_two_id": filter.Person2})
	}

	query := stBuilder.
		Select("*").
		From(LikeTableName).
		Where(whereMap).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Like, 0)
	var tmp models.Like
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db can't scan: ", err.Error())
			return nil, err
		}
		res = append(res, &tmp)
	}

	return res, nil
}

func (storage *LikeStorage) Create(requestID int64, person1ID, person2ID types.UserID) error {
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db create request to ", LikeTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Insert(LikeTableName).
		Columns("person_one_id", "person_two_id").
		Values(person1ID, person2ID).
		RunWith(storage.dbReader)

	_, err := query.Query()
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't create like: ", err.Error())
	}
	return err
}

func (storage *LikeStorage) GetMatch(requestID int64, person1ID types.UserID) ([]types.UserID, error) {
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db get request to ", LikeTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.Select("t1.person_two_id").
		From(LikeTableName + "t1").
		InnerJoin(LikeTableName + "t2 ON t1.person_one_id = t2.person_two_id AND t1.person_two_id = t2.person_one_id").
		Where(qb.Eq{"t1.person_one_id": person1ID}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't query:  ", err.Error())
		return nil, err
	}

	res := make([]types.UserID, 0)
	var scan types.UserID
	for rows.Next() {
		err = rows.Scan(&scan)
		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't scan:  ", err.Error())
			return nil, err
		}
		res = append(res, scan)
	}
	return res, nil
}
