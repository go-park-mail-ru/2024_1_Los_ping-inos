package repo

import (
	"context"
	"database/sql"

	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	"main.go/internal/feed"
	. "main.go/internal/logs"
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

func (storage *LikeStorage) Get(ctx context.Context, filter *feed.LikeGetFilter) ([]types.UserID, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to ", LikeTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	whereMap := qb.Or{}

	if filter == nil {
		filter = &feed.LikeGetFilter{}
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
	var tmp feed.Like
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
	logger := ctx.Value(Logg).(Log)
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

func (storage *LikeStorage) GetLikesLeft(ctx context.Context, personID types.UserID) (int, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db create request to ", LikeTableName)

	query := "SELECT likes_left FROM person WHERE id = $1"

	rows, err := storage.dbReader.Query(query, personID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var likes int

	err = rows.Scan(likes)
	if err != nil {
		return 0, err
	}

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("halliluya")
	return likes, nil
}

func (storage *LikeStorage) DecreaseLikesCount(ctx context.Context, personID types.UserID) (int, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db create request to ", LikeTableName)

	query := "SELECT likes_left FROM person WHERE id = $1"

	rows, err := storage.dbReader.Query(query, personID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var likes int

	for rows.Next() {
		err = rows.Scan(&likes)
		if err != nil {
			return 0, err
		}
	}

	likes--

	_, err = storage.dbReader.Exec("UPDATE person SET likes_left = $1 WHERE id = $2", likes, personID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return 0, err
	}

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("halliluya")
	return likes, nil
}

func (storage *LikeStorage) IncreaseLikesCount(ctx context.Context, personID types.UserID) error {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db create request to ", LikeTableName)

	_, err := storage.dbReader.Exec("UPDATE person SET likes_left = $1 WHERE id = $2", 5, personID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return err
	}

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("halliluya")
	return nil
}
