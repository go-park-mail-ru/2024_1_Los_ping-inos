package storage

import (
	"database/sql"
	"fmt"

	qb "github.com/Masterminds/squirrel"
	models "main.go/db"
)

type Storage struct {
	Person PersonStorage
}

type PersonStorage struct {
	dbReader *sql.DB
}

func (storage *PersonStorage) Get(filter *models.PersonFilter) ([]*models.Person, error) {
	whereMap := qb.And{}

	if filter == nil {
		filter = &models.PersonFilter{}
	}

	whereMap = processIDFilter(filter, whereMap)
	whereMap = processEmailFilter(filter, whereMap)
	whereMap = processSessionIDFilter(filter, whereMap)

	query := qb.
		Select("*").    // select * - not safe
		From("person"). // TODO название таблиц в константы
		Where(whereMap).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons := make([]*models.Person, 0)

	for rows.Next() {
		person := &models.Person{}
		err := rows.Scan(&person.ID, &person.Name, &person.Birthday, &person.Description, &person.Location, &person.Photo,
			&person.Email, &person.Password, &person.CreatedAt, &person.Premium, &person.LikesLeft, &person.SessionID)

		if err != nil {
			return nil, err
		}

		persons = append(persons, person)
	}

	return persons, nil
}

func (storage *PersonStorage) AddAccount(Name string, Birthday string, Gender string, Email string, Password string) error {
	_, err := storage.dbReader.Exec(
		"INSERT INTO person(name, birthday, email, password, gender) "+
			"VALUES (s1, s2, s3, s4, s5)", Name, Birthday, Email, Password, Gender)
	if err != nil {
		return fmt.Errorf("Create user %w", err)
	}

	return nil
}

func processIDFilter(filter *models.PersonFilter, whereMap qb.And) qb.And { // what is it?
	if filter.ID != nil {
		IDMap := qb.Or{}
		for _, id := range filter.ID {
			IDMap = append(IDMap, qb.Eq{"id": id})
		}
		whereMap = append(whereMap, IDMap)
	}
	return whereMap
}

func processEmailFilter(filter *models.PersonFilter, whereMap qb.And) qb.And {
	if filter.Email != nil {
		emailMap := qb.Or{}
		for _, ID := range filter.Email {
			emailMap = append(emailMap, qb.Eq{"email": ID})
		}
		whereMap = append(whereMap, emailMap)
	}
	return whereMap
}

func processSessionIDFilter(filter *models.PersonFilter, whereMap qb.And) qb.And {
	if filter.SessionID != nil {
		emailMap := qb.Or{}
		for _, ID := range filter.Email {
			emailMap = append(emailMap, qb.Eq{"session_id": ID})
		}
		whereMap = append(whereMap, emailMap)
	}
	return whereMap
}
