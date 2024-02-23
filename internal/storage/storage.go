package storage

import (
	"database/sql"
	"fmt"

	"main.go/internal/types"
)

type Storage struct {
	db *sql.DB
}

// Get - вообще он должен возвращать какую-то структуру, которая лежит в models, но для примера пока так
func (e Storage) Get(id types.UserID) string {
	return fmt.Sprintf("cool id %v!\n", id)
}
