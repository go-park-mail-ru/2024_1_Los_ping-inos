package repo

import (
	"context"
	"sync"
	"time"

	qb "github.com/Masterminds/squirrel"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"main.go/internal/feed"
	. "main.go/internal/logs"
	"main.go/internal/types"
)

const (
	messageFields    = "id, data, sender_id, receiver_id, sent_time"
	messageTable     = "message"
	getLastMessQuery = "SELECT DISTINCT ON (    CASE        WHEN sender_id < receiver_id THEN sender_id || '_' || receiver_id        ELSE receiver_id || '_' || sender_id    END) id, data, sender_id, receiver_id, sent_time FROM message WHERE     (sender_id = $1 OR receiver_id = $1)    AND ((sender_id = ANY($2)) OR (receiver_id = ANY($2))) ORDER BY (    CASE        WHEN sender_id < receiver_id THEN sender_id || '_' || receiver_id        ELSE receiver_id || '_' || sender_id    END), sent_time DESC;"
)

func (storage *PostgresStorage) GetChat(ctx context.Context, user1, user2 types.UserID) ([]feed.Message, string, error) {
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
		return nil, "", err
	}
	defer rows.Close()

	var (
		messages []feed.Message
		message  feed.Message
	)

	for rows.Next() {
		err = rows.Scan(&message.Id, &message.Data, &message.Sender, &message.Receiver, &message.Time)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn("Db can't scan: ", err.Error())
			return nil, "", err
		}
		messages = append(messages, message)
	}

	qq := "SELECT name FROM person WHERE id = $1"

	rows, err = storage.dbReader.Query(qq, user2)
	if err != nil {
		return []feed.Message{}, "", err
	}
	defer rows.Close()

	var name string

	for rows.Next() {

		err := rows.Scan(name)
		if err != nil {
			return []feed.Message{}, "", err
		}
	}

	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Db: send messages")
	return messages, name, nil
}

func (storage *PostgresStorage) CreateMessage(ctx context.Context, message feed.MessageToReceive) (*feed.MessageToReceive, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Create request to message")
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	t := time.UnixMilli(message.Time)
	query := stBuilder.
		Insert(messageTable).
		Columns("data, sender_id, receiver_id, sent_time").
		Values(message.Data, message.Sender, message.Receiver, t).
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

func (storage *PostgresStorage) GetLastMessages(ctx context.Context, id int64, ids []int) ([]feed.Message, error) {
	logger := ctx.Value(Logg).(Log)
	logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("Get last messages request")
	rows, err := storage.dbReader.Query(getLastMessQuery, id, pq.Array(ids))
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("db can't query: ", err.Error())
		return nil, err
	}
	defer rows.Close()
	var res []feed.Message
	for rows.Next() {
		tmp := feed.Message{}
		rows.Scan(&tmp.Id, &tmp.Data, &tmp.Sender, &tmp.Receiver, &tmp.Time)
		res = append(res, tmp)
	}
	return res, nil
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
