package repo

import (
	"context"
	qb "github.com/Masterminds/squirrel"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"main.go/internal/feed"
	. "main.go/internal/logs"
	"main.go/internal/types"
	"sync"
)

const (
	messageFields = "data, sender_id, receiver_id, sent_time"
	messageTable  = "message"
)

func (storage *PostgresStorage) GetChat(ctx context.Context, user1, user2 types.UserID) ([]feed.Message, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Get request to message")

	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Select(messageFields).
		From(messageTable).
		Where(qb.Or{qb.And{qb.Eq{"sender_id": user1}, qb.Eq{"receiver_id": user2}},
			qb.And{qb.Eq{"sender_id": user2}, qb.Eq{"receiver_id": user1}}}).
		OrderBy("sent_time").
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Db can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	var (
		messages []feed.Message
		message  feed.Message
	)

	for rows.Next() {
		err = rows.Scan(&message.Data, &message.Sender, &message.Receiver, &message.Time)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Db can't scan: ", err.Error())
			return nil, err
		}
		messages = append(messages, message)
	}

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Db: send messages")
	return messages, nil
}

func (storage *PostgresStorage) CreateMessage(ctx context.Context, message feed.Message) (*feed.Message, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Create request to message")

	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Insert(messageTable).
		Columns(messageFields).
		Values(message.Data, message.Sender, message.Receiver, message.Time).
		RunWith(storage.dbReader)

	rows, err := query.Query()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Db can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Db: created message")
	return &message, nil
}

type WSStorage struct {
	connections *sync.Map // [UID]*websocket.Connection
}

func NewWebSocStorage() *WSStorage {
	return &WSStorage{
		connections: &sync.Map{},
	}
}

func (stor *WSStorage) AddConnection(ctx context.Context, connection *websocket.Conn, UID types.UserID) {
	stor.connections.Store(UID, connection)
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Stored connection with user ", UID)
}

func (stor *WSStorage) GetConnection(ctx context.Context, UID types.UserID) (*websocket.Conn, bool) {
	logger := ctx.Value(Logg).(Log)
	con, ok := stor.connections.Load(UID)
	if ok {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Got no connection with user ", UID)
		return con.(*websocket.Conn), ok
	}
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Got connection with user ", UID)
	return nil, ok
}

func (stor *WSStorage) DeleteConnection(ctx context.Context, UID types.UserID) {
	stor.connections.Delete(UID)
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Deleted connection with user ", UID)
}