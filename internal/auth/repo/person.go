package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"main.go/internal/auth"
	. "main.go/internal/logs"
	"main.go/internal/types"
)

const (
	personFields    = "id, name, birthday, description, location, photo, email, password, created_at, premium, likes_left, session_id, gender"
	PersonTableName = "person"
	LikeTableName   = "\"like\""
)

type PersonStorage struct {
	dbReader *sql.DB
}

func NewAuthPersonStorage(dbReader *sql.DB) *PersonStorage {
	return &PersonStorage{
		dbReader: dbReader,
	}
}

func (storage *PersonStorage) Get(ctx context.Context, filter *auth.PersonGetFilter) ([]*auth.Person, error) {
	logger := ctx.Value(Logg).(Log)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	whereMap := qb.And{}

	if filter == nil {
		filter = &auth.PersonGetFilter{}
	}

	processIDFilter(filter, &whereMap)
	processEmailFilter(filter, &whereMap)
	processSessionIDFilter(filter, &whereMap)

	query := stBuilder.
		Select(personFields).
		From(PersonTableName).
		Where(whereMap).
		RunWith(storage.dbReader)

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to ", PersonTableName)
	rows, err := query.Query()

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db read can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	persons := make([]*auth.Person, 0)

	for rows.Next() {
		person := &auth.Person{}
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

func (storage *PersonStorage) Update(ctx context.Context, person auth.Person) error {
	logger := ctx.Value(Logg).(Log)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	setMap := make(map[string]interface{})
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db update request to ", PersonTableName)

	tmp, err := json.Marshal(person)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
		return err
	}
	err = json.Unmarshal(tmp, &setMap)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
		return err
	}
	setMap["password"] = person.Password
	query := stBuilder.
		Update(PersonTableName).
		SetMap(setMap).
		Where(qb.Eq{"id": person.ID}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	defer rows.Close()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db update can't query: ", err.Error())
		return err
	}
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db person updated")
	return nil
}

func (storage *PersonStorage) Delete(ctx context.Context, sessionID string) error {
	logger := ctx.Value(Logg).(Log)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db delete request to ", PersonTableName)
	query := stBuilder.
		Delete(PersonTableName).
		Where(qb.Eq{"session_id": sessionID}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	defer rows.Close()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db delete can't query: ", err.Error())
		return err
	}
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db person deleted")
	return nil
}

func (storage *PersonStorage) AddAccount(ctx context.Context, Name string, Birthday string, Gender string, Email string, Password string) error {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db create request to ", PersonTableName)
	_, err := storage.dbReader.Exec(
		"INSERT INTO person(name, birthday, email, password, gender) "+
			"VALUES ($1, $2, $3, $4, $5)", Name, Birthday, Email, Password, Gender)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't query: ", err.Error())

		return err
	}

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db created person")
	return nil
}

func (storage *PersonStorage) RemoveSession(ctx context.Context, sid string) error {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db remove session_id request to ", PersonTableName)
	_, err := storage.dbReader.Exec(
		"UPDATE person SET session_id = '' "+
			"WHERE session_id = $1", sid)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Remove sessions %w", err)
	}
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db removed session_id ", PersonTableName)
	return nil
}

func (storage *PersonStorage) GetMatch(ctx context.Context, person1ID types.UserID) ([]types.UserID, error) {
	logger := ctx.Value(Logg).(Log)
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

func processIDFilter(filter *auth.PersonGetFilter, whereMap *qb.And) {
	if filter.ID != nil {
		*whereMap = append(*whereMap, qb.Eq{"id": filter.ID})
	}
}

func processEmailFilter(filter *auth.PersonGetFilter, whereMap *qb.And) {
	if filter.Email != nil {
		*whereMap = append(*whereMap, qb.Eq{"email": filter.Email})
	}
}

func processSessionIDFilter(filter *auth.PersonGetFilter, whereMap *qb.And) {
	if filter.SessionID != nil {
		*whereMap = append(*whereMap, qb.Eq{"session_id": filter.SessionID})
	}
}
