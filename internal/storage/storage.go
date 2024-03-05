package storage

import (
	"database/sql"
	"fmt"

	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"main.go/db"
)

type PersonStorage struct {
	dbReader *sql.DB
}

func NewPersonStorage(dbReader *sql.DB) *PersonStorage {
	return &PersonStorage{
		dbReader: dbReader,
	}
}

func (storage *PersonStorage) Get(filter *models.PersonFilter) ([]*models.Person, error) {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	whereMap := qb.And{}

	if filter == nil {
		filter = &models.PersonFilter{}
	}

	processIDFilter(filter, &whereMap)
	processEmailFilter(filter, &whereMap)
	processSessionIDFilter(filter, &whereMap)

	query := stBuilder.
		Select("*").
		From("person"). // TODO название таблиц в константы
		Where(whereMap).
		RunWith(storage.dbReader)

	rows, err := query.Query()

	if err != nil {
		logrus.Info("can't query")
		println(query.ToSql())
		println(err.Error())
		return nil, err
	}
	defer rows.Close()

	persons := make([]*models.Person, 0)

	for rows.Next() {
		person := &models.Person{}
		err := rows.Scan(&person.ID, &person.Name, &person.Birthday, &person.Description, &person.Location, &person.Photo,
			&person.Email, &person.Password, &person.Gender, &person.CreatedAt, &person.Premium, &person.LikesLeft, &person.SessionID)

		if err != nil {
			logrus.Info("can't scan row ", err.Error())
			return nil, err
		}

		persons = append(persons, person)
	}

	return persons, nil
}

func (storage *PersonStorage) AddAccount(Name string, Birthday string, Gender string, Email string, Password string) error {
	_, err := storage.dbReader.Exec(
		"INSERT INTO person(name, birthday, email, password, gender) "+
			"VALUES ($1, $2, $3, $4, $5)", Name, Birthday, Email, Password, Gender)
	if err != nil {
		println(err.Error())
		return fmt.Errorf("Create user %w", err)
	}

	return nil
}

func processIDFilter(filter *models.PersonFilter, whereMap *qb.And) {
	if filter.ID != nil {
		*whereMap = append(*whereMap, qb.Eq{"id": filter.ID})
	}
}

func processEmailFilter(filter *models.PersonFilter, whereMap *qb.And) {
	if filter.Email != nil {
		*whereMap = append(*whereMap, qb.Eq{"email": filter.Email})
	}

}

func processSessionIDFilter(filter *models.PersonFilter, whereMap *qb.And) {
	if filter.SessionID != nil {
		*whereMap = append(*whereMap, qb.Eq{"session_id": filter.SessionID})
	}
}
