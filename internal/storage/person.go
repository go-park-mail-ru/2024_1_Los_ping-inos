package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	. "main.go/config"
	models "main.go/db"
	. "main.go/internal/logs"
)

type PersonStorage struct {
	dbReader *sql.DB
}

func NewPersonStorage(dbReader *sql.DB) *PersonStorage {
	return &PersonStorage{
		dbReader: dbReader,
	}
}

func (storage *PersonStorage) Get(requestID int64, filter *models.PersonGetFilter) ([]*models.Person, error) {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	whereMap := qb.And{}

	if filter == nil {
		filter = &models.PersonGetFilter{}
	}

	processIDFilter(filter, &whereMap)
	processEmailFilter(filter, &whereMap)
	processSessionIDFilter(filter, &whereMap)

	query := stBuilder.
		Select("*").
		From(PersonTableName).
		Where(whereMap).
		RunWith(storage.dbReader)

	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db get request to ", PersonTableName)
	rows, err := query.Query()

	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	persons := make([]*models.Person, 0)

	for rows.Next() {
		person := &models.Person{}
		err = rows.Scan(&person.ID, &person.Name, &person.Birthday, &person.Description, &person.Location, &person.Photo,
			&person.Email, &person.Password, &person.CreatedAt, &person.Premium, &person.LikesLeft, &person.SessionID, &person.Gender)

		if err != nil {
			Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't scan person: ", err.Error())
			return nil, err
		}

		persons = append(persons, person)
	}

	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db returning records")
	return persons, nil
}

func (storage *PersonStorage) Update(requestID int64, person models.Person) error {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	setMap := make(map[string]interface{})
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db update request to ", PersonTableName)

	tmp, err := json.Marshal(person)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
		return err
	}
	err = json.Unmarshal(tmp, &setMap)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn(err.Error())
		return err
	}

	query := stBuilder.
		Update(PersonTableName).
		SetMap(setMap).
		Where(qb.Eq{"id": person.ID}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	defer rows.Close()
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't query: ", err.Error())
		return err
	}
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db person updated")
	return nil
}

func (storage *PersonStorage) AddAccount(requestID int64, Name string, Birthday string, Gender string, Email string, Password string) error {
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db create request to ", PersonTableName)
	_, err := storage.dbReader.Exec(
		"INSERT INTO person(name, birthday, email, password, gender) "+ // TODO PersonTableName
			"VALUES ($1, $2, $3, $4, $5)", Name, Birthday, Email, Password, Gender)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("db can't query: ", err.Error())

		return err
	}

	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db created person")
	return nil
}

func (storage *PersonStorage) RemoveSession(requestID int64, sid string) error {
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db remove session_id request to ", PersonTableName)
	_, err := storage.dbReader.Exec(
		"UPDATE person SET session_id = '' "+
			"WHERE session_id = $1", sid)
	if err != nil {
		Log.WithFields(logrus.Fields{RequestID: requestID}).Warn("can't query: ", err.Error())
		return fmt.Errorf("Remove sessions %w", err)
	}
	Log.WithFields(logrus.Fields{RequestID: requestID}).Info("db removed session_id ", PersonTableName)
	return nil
}

func processIDFilter(filter *models.PersonGetFilter, whereMap *qb.And) {
	if filter.ID != nil {
		*whereMap = append(*whereMap, qb.Eq{"id": filter.ID})
	}
}

func processEmailFilter(filter *models.PersonGetFilter, whereMap *qb.And) {
	if filter.Email != nil {
		*whereMap = append(*whereMap, qb.Eq{"email": filter.Email})
	}

}

func processSessionIDFilter(filter *models.PersonGetFilter, whereMap *qb.And) {
	if filter.SessionID != nil {
		*whereMap = append(*whereMap, qb.Eq{"session_id": filter.SessionID})
	}
}
