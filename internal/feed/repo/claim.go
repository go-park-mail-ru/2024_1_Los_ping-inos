package repo

import (
	"context"
	qb "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"main.go/internal/feed"
	. "main.go/internal/logs"
)

const (
	claimFields = "type, sender_id, receiver_id"
)

func (storage *PostgresStorage) CreateClaim(ctx context.Context, claim feed.Claim) error {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db create request to person_claim")
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Insert("person_claim").
		Columns(claimFields).
		Values(claim.TypeID, claim.SenderID, claim.ReceiverID).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't create claim: ", err.Error())
		return err
	}
	defer rows.Close()
	return nil
}
