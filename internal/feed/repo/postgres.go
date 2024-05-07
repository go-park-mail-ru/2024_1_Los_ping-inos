package repo

import "database/sql"

type PostgresStorage struct {
	dbReader *sql.DB
}

func NewPostgresStorage(dbReader *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		dbReader: dbReader,
	}
}
