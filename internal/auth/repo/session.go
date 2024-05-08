package repo

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"main.go/internal/auth"
	. "main.go/internal/logs"
	"main.go/internal/types"
)

type SessionStorage struct {
	db *redis.Client
}

func NewSessionStorage(db *redis.Client) *SessionStorage {
	return &SessionStorage{
		db: db,
	}
}

func (stor *SessionStorage) GetBySID(ctx context.Context, SID string) (*auth.Session, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db get request to session storage")
	session := &auth.Session{SID: SID}

	res, err := stor.db.Get(context.TODO(), SID).Result()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't query: ", err.Error())
		return nil, err
	}

	UID, err := strconv.Atoi(res)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("db can't query: ", err.Error())
		return nil, err
	}
	session.UID = types.UserID(UID)
	return session, nil
}

func (stor *SessionStorage) CreateSession(ctx context.Context, UID types.UserID) (string, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db create request to session storage")
	SID := uuid.NewString()
	err := stor.db.Set(context.TODO(), SID, strconv.Itoa(int(UID)), 0).Err()
	return SID, err
}

func (stor *SessionStorage) DeleteSession(ctx context.Context, SID string) error {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db delete request to session storage")
	return stor.db.Del(context.TODO(), SID).Err()
}
