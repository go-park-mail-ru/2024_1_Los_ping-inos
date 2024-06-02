package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	models "main.go/internal/feed"
	mocks "main.go/internal/feed/mocks"
	"main.go/internal/types"
)

func TestGetLastMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []models.Message{
		{
			MsgType: "message",
			Properties: models.MsgProperties{
				Id:       1,
				Data:     "just got a new glock, its fif'een",
				Sender:   1,
				Receiver: 2,
				Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			MsgType: "message",
			Properties: models.MsgProperties{
				Id:       2,
				Data:     "just got a new glock, its fif'een",
				Sender:   1,
				Receiver: 2,
				Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	mockObj := mocks.NewMockPostgresStorage(ctrl)
	mockObj.EXPECT().GetLastMessages(gomock.Any(), int64(1), []int{1, 2, 3}).Return(expected, nil)

	core := UseCase{storage: mockObj}

	claim := []models.Message{
		{
			MsgType: "message",
			Properties: models.MsgProperties{
				Id:       1,
				Data:     "just got a new glock, its fif'een",
				Sender:   1,
				Receiver: 2,
				Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			MsgType: "message",
			Properties: models.MsgProperties{
				Id:       2,
				Data:     "just got a new glock, its fif'een",
				Sender:   1,
				Receiver: 2,
				Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	result, err := core.GetLastMessages(context.TODO(), int64(1), []int64{1, 2, 3})
	if err != nil {
		t.Errorf("unexpected err result")
		t.Error(err)
		return
	}
	require.Equal(t, claim, result)
}

func TestAddConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var nilConnection *websocket.Conn
	newUID := []types.UserID{
		123,
		231,
	}

	mockObj := mocks.NewMockWebSocStorage(ctrl)
	mockObj.EXPECT().AddConnection(gomock.Any(), nilConnection, newUID[0]).Return(nil)
	mockObj.EXPECT().AddConnection(gomock.Any(), nilConnection, newUID[1]).Return(fmt.Errorf("repo error"))

	service := UseCase{ws: mockObj}

	err := service.AddConnection(context.TODO(), nilConnection, newUID[0])
	if err != nil {
		t.Errorf("unexpected err result")
		t.Error(err)
		return
	}
	err = service.AddConnection(context.TODO(), nilConnection, newUID[1])
	if err == nil {
		t.Errorf("unexpected err result")
		t.Error(err)
		return
	}
}

func TestGetConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var nilConnection *websocket.Conn
	newUID := []types.UserID{
		123,
		231,
	}

	mockObj := mocks.NewMockWebSocStorage(ctrl)
	mockObj.EXPECT().GetConnection(gomock.Any(), newUID[0]).Return(nilConnection, true)
	//mockObj.EXPECT().AddConnection(gomock.Any(), nilConnection, newUID[1]).Return(fmt.Errorf("repo error"))

	service := UseCase{ws: mockObj}

	conn, ok := service.GetConnection(context.TODO(), newUID[0])
	if conn != nilConnection {
		t.Errorf("unexpected err result")
		t.Error(conn)
		return
	}
	if ok != true {
		t.Errorf("unexpected err result")
		t.Error(ok)
		return
	}

}

func TestDeletetConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newUID := []types.UserID{
		123,
		231,
	}

	mockObj := mocks.NewMockWebSocStorage(ctrl)
	mockObj.EXPECT().DeleteConnection(gomock.Any(), newUID[0]).Return(nil)
	mockObj.EXPECT().DeleteConnection(gomock.Any(), newUID[1]).Return(fmt.Errorf("repo error"))
	//mockObj.EXPECT().AddConnection(gomock.Any(), nilConnection, newUID[1]).Return(fmt.Errorf("repo error"))

	service := UseCase{ws: mockObj}

	err := service.DeleteConnection(context.TODO(), newUID[0])
	if err != nil {
		t.Errorf("unexpected err result")
		t.Error(err)
		return
	}
	err = service.DeleteConnection(context.TODO(), newUID[1])
	if err == nil {
		t.Errorf("unexpected err result")
		t.Error(err)
		return
	}

}
