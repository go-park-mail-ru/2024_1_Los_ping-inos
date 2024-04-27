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
	csatColumns   = "q1, tittle_id"
)

type CsatStorage struct {
	dbReader *sql.DB
}

func NewCsatStorage(dbReader *sql.DB) *CsatStorage {
	return &CsatStorage{
		dbReader: dbReader,
	}
}

func (storage CsatStorage) Create(ctx context.Context, q1 int, tittleID int) error {
	logger := ctx.Value(Logg).(Log)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to ", csatTableName)

	query := stBuilder.
		Insert(csatTableName).
		Columns(csatColumns).
		Values(q1, tittleID).
		RunWith(storage.dbReader)

	rows, err := query.Query()

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't create rate: ", err.Error())
		return err
	}
	defer rows.Close()
	return nil
}

func (storage CsatStorage) GetTittlesCount(ctx context.Context) (int, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db create request to ")

	var count int

	err := storage.dbReader.QueryRow("SELECT COUNT(DISTINCT tittle) FROM questions").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (storage CsatStorage) GetStat(ctx context.Context, tittleID int) (string, float32, []int, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db create request to ")

	//query := "SELECT * FROM "
	var tittle string

	err := storage.dbReader.QueryRow("SELECT tittle from questions WHERE id = $1", tittleID).Scan(&tittle)
	if err != nil {
		return "", 0, nil, err
	}

	var stats []int

	query := "SELECT q1 FROM csat WHERE tittle_id = $1"

	rows, err := storage.dbReader.Query(query, tittleID)
	if err != nil {
		return "", 0, nil, err
	}

	for rows.Next() {
		var q int
		err = rows.Scan(&q)
		if err != nil {
			return "", 0, nil, err
		}
		stats = append(stats, q)
	}

	var count, summary int
	for _, number := range stats {
		summary += number
		count++
	}

	//println(summary, count)
	avg := float32(summary) / float32(count)
	//println(avg)

	//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Return  images")
	return tittle, avg, stats, nil
}
