package repo

import (
	"context"
	"database/sql"
	"main.go/internal/types"

	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	"main.go/internal/feed"
	. "main.go/internal/logs"
)

const (
	personFields = "id, name, birthday, description, location, photo, email, password, created_at, premium, likes_left, session_id, gender"
)

type PersonStorage struct {
	dbReader *sql.DB
}

func NewPersonStorage(dbReader *sql.DB) *PersonStorage {
	return &PersonStorage{
		dbReader: dbReader,
	}
}

func (storage *PersonStorage) GetFeed(ctx context.Context, filter types.UserID) ([]*feed.Person, error) {
	logger := ctx.Value(Logg).(Log)
	likes := &LikeStorage{dbReader: storage.dbReader}
	ids, err := likes.Get(ctx, &feed.LikeGetFilter{Person1: &filter})
	if err != nil {
		return nil, err
	}
	ids = append(ids, filter)

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
		err = rows.Scan(&person.ID, &person.Name, &person.Birthday, &person.Description, &person.Location, &person.Photo,
			&person.Email, &person.Password, &person.CreatedAt, &person.Premium, &person.LikesLeft, &person.SessionID, &person.Gender)

		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't scan person: ", err.Error())
			return nil, err
		}

		persons = append(persons, person)
	}

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db returning records")
	return persons, nil
}
