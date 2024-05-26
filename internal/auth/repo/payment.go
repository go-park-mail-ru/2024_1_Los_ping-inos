package repo

import (
	"context"
	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"main.go/internal/auth"
	. "main.go/internal/logs"
	"main.go/internal/types"
	"time"
)

func (storage *PersonStorage) ActivateSub(ctx context.Context, UID types.UserID, datetime time.Time) error {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db activating sub ", PersonTableName)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	setMap := map[string]interface{}{}
	setMap["premium"] = true
	setMap["premium_expires_at"] = datetime.Add(31 * 24 * time.Hour)
	query := stBuilder.Update(PersonTableName).SetMap(setMap).Where(qb.Eq{"id": UID}).RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't add premium: ", err.Error())
		return err
	}
	defer rows.Close()
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("premium activated for user ", UID)

	query2 := stBuilder.
		Insert("person_payment").
		Columns("person_id", "paymentTime").
		Values(UID, datetime).
		RunWith(storage.dbReader)
	rows2, err := query2.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't add premium to history: ", err.Error())
		return err
	}
	defer rows2.Close()
	return nil
}

func (storage *PersonStorage) GetSubHistory(ctx context.Context, UID types.UserID) (*auth.PaymentHistory, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db getting sub history for user ", UID)
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Select("paymentTime").
		From("person_payment").
		Where(qb.Eq{"person_id": UID}).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't read sub history: ", err.Error())
		return nil, err
	}

	res := &auth.PaymentHistory{}
	var tmp time.Time
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't scan sub history row: ", err.Error())
			return nil, err
		}
		res.Times = append(res.Times, tmp.Unix())
	}

	return res, nil
}
