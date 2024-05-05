package feed

import (
	"context"
	"github.com/gorilla/websocket"
	"main.go/internal/types"
)

type (
	UseCase interface {
		GetCards(userID types.UserID, ctx context.Context) ([]Card, error)
		CreateLike(profile1, profile2 types.UserID, ctx context.Context) error
		GetChat(ctx context.Context, user1, user2 types.UserID) ([]Message, error)
		GetLastMessages(ctx context.Context, UID int64, ids []int64) ([]Message, error)

		AddConnection(ctx context.Context, connection *websocket.Conn, UID types.UserID)
		GetConnection(ctx context.Context, UID types.UserID) (*websocket.Conn, bool)
		DeleteConnection(ctx context.Context, UID types.UserID)
		SaveMessage(ctx context.Context, message Message) (*Message, error)
		CreateClaim(ctx context.Context, typeID, senderID, receiverID int64) error
		GetClaims(ctx context.Context) ([]PureClaim, error)
	}
	PostgresStorage interface {
		GetFeed(ctx context.Context, filter types.UserID) ([]*Person, error)
		GetPersonInterests(ctx context.Context, personID types.UserID) ([]*Interest, error)
		GetLike(ctx context.Context, filter *LikeGetFilter) ([]types.UserID, error)
		CreateLike(ctx context.Context, person1ID, person2ID types.UserID) error
		GetImages(ctx context.Context, userID int64) ([]Image, error)
		GetChat(ctx context.Context, user1, user2 types.UserID) ([]Message, error)
		CreateMessage(ctx context.Context, message Message) (*Message, error)
		GetLastMessages(ctx context.Context, id int64, ids []int) ([]Message, error)
		CreateClaim(ctx context.Context, claim Claim) error
		GetAllClaims(ctx context.Context) ([]PureClaim, error)
	}

	WebSocStorage interface {
		AddConnection(ctx context.Context, connection *websocket.Conn, UID types.UserID)
		GetConnection(ctx context.Context, UID types.UserID) (*websocket.Conn, bool)
		DeleteConnection(ctx context.Context, UID types.UserID)
	}
)
