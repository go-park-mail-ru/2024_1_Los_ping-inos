package repo

import (
	"context"
	"main.go/internal/types"

	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	"main.go/internal/feed"
	. "main.go/internal/logs"
)

const (
	personFields = "id, name, birthday, description, location, email, password, created_at, premium, likes_left, gender"
)

func (storage *PostgresStorage) GetFeed(ctx context.Context, filter types.UserID) ([]*feed.Person, error) {
	logger := ctx.Value(Logg).(Log)

	ids, err := storage.GetLike(ctx, &feed.LikeGetFilter{Person1: &filter})
	if err != nil {
		return nil, err
	}
	ids = append(ids, filter)

	baned, err := storage.GetClaimed(ctx, filter)
	if err != nil {
		return nil, err
	}
	ids = append(ids, baned...)

	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to ", PersonTableName)
	query := stBuilder.
		Select(personFields).
		From(PersonTableName).
		Where(qb.NotEq{"id": ids}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db read can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	persons := make([]*feed.Person, 0)

	for rows.Next() {
		person := &feed.Person{}
		err = rows.Scan(&person.ID, &person.Name, &person.Birthday, &person.Description, &person.Location,
			&person.Email, &person.Password, &person.CreatedAt, &person.Premium, &person.LikesLeft, &person.Gender)

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't scan person: ", err.Error())
			return nil, err
		}

		persons = append(persons, person)
	}

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db returning records")
	return persons, nil
}

func (storage *PostgresStorage) GetClaimed(ctx context.Context, id types.UserID) ([]types.UserID, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to person_claim")
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Select("receiver_id").
		From("person_claim").Where(qb.Eq{"sender_id": id}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()
	var res []types.UserID
	var tmp types.UserID
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't scan: ", err.Error())
			return nil, err
		}
		res = append(res, tmp)
	}
	return res, nil
}
