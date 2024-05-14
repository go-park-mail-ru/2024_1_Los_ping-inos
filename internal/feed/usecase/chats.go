package usecase

import (
	"context"

	"github.com/gorilla/websocket"
	"main.go/internal/feed"
	"main.go/internal/types"
)

func (service *UseCase) AddConnection(ctx context.Context, connection *websocket.Conn, UID types.UserID) error {
	err := service.ws.AddConnection(ctx, connection, UID)
	if err != nil {
		return err
	}
	return nil
}
func (service *UseCase) GetConnection(ctx context.Context, UID types.UserID) (*websocket.Conn, bool) {
	return service.ws.GetConnection(ctx, UID)
}

func (service *UseCase) DeleteConnection(ctx context.Context, UID types.UserID) error {
	err := service.ws.DeleteConnection(ctx, UID)
	if err != nil {
		return err
	}
	return nil
}

func (service *UseCase) GetLastMessages(ctx context.Context, UID int64, ids []int64) ([]feed.Message, error) {
	input := make([]int, len(ids))
	for i := range ids {
		input[i] = int(ids[i])
	}
	return service.storage.GetLastMessages(ctx, UID, input)
}
