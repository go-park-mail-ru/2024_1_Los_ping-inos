package storage

import (
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"main.go/db"
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

	if filter.ID != nil {
		IDMap := qb.Or{}
		for _, id := range filter.ID {
			IDMap = append(IDMap, qb.Eq{"id": id})
		}
		whereMap = append(whereMap, IDMap)
	}

	if filter.Email != nil {
		emailMap := qb.Or{}
		for _, ID := range filter.Email {
			emailMap = append(emailMap, qb.Eq{"email": ID})
		}
		whereMap = append(whereMap, emailMap)
	}

	if filter.SessionID != nil {
		emailMap := qb.Or{}
		for _, ID := range filter.Email {
			emailMap = append(emailMap, qb.Eq{"session_id": ID})
		}
		whereMap = append(whereMap, emailMap)
	}

	query := qb.
		Select("*").
		From("person").
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
