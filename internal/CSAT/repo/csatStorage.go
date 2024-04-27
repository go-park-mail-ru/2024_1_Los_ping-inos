package repo

import (
	"context"
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	. "main.go/internal/logs"
)

const (
	csatTableName = "csat"
	csatColumns   = "q1"
)

type CsatStorage struct {
	dbReader *sql.DB
}

func NewCsatStorage(dbReader *sql.DB) *CsatStorage {
	return &CsatStorage{
		dbReader: dbReader,
	}
}

func (storage CsatStorage) Create(ctx context.Context, q1 int) error {
	logger := ctx.Value(Logg).(Log)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to ", csatTableName)

	query := stBuilder.
		Insert(csatTableName).
		Columns(csatColumns).
		Values(q1).
		RunWith(storage.dbReader)

	rows, err := query.Query()

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't create rate: ", err.Error())
		return err
	}
	defer rows.Close()
	return nil
}

func (storage CsatStorage) GetStat(ctx context.Context) (map[string]int, error) {
	return nil, nil
}
