package usecase

import (
	"context"
	"github.com/gorilla/websocket"
	"main.go/internal/types"
)

func (service *UseCase) AddConnection(ctx context.Context, connection *websocket.Conn, UID types.UserID) {
	service.ws.AddConnection(ctx, connection, UID)
}
func (service *UseCase) GetConnection(ctx context.Context, UID types.UserID) (*websocket.Conn, bool) {
	return service.ws.GetConnection(ctx, UID)
}

func (service *UseCase) DeleteConnection(ctx context.Context, UID types.UserID) {
	service.ws.DeleteConnection(ctx, UID)
}
