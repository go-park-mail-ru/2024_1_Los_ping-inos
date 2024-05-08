package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	models "main.go/internal/feed"
	mocks "main.go/internal/feed/mocks"
)

func TestGetLastMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []models.Message{
		{
			Id:       1,
			Data:     "just got a new glock, its fif'een",
			Sender:   1,
			Receiver: 2,
			Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:       2,
			Data:     "just got a new glock, its fif'een",
			Sender:   1,
			Receiver: 2,
			Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	mockObj := mocks.NewMockPostgresStorage(ctrl)
	mockObj.EXPECT().GetLastMessages(gomock.Any(), int64(1), []int{1, 2, 3}).Return(expected, nil)

	core := UseCase{storage: mockObj}

	claim := []models.Message{
		{
			Id:       1,
			Data:     "just got a new glock, its fif'een",
			Sender:   1,
			Receiver: 2,
			Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:       2,
			Data:     "just got a new glock, its fif'een",
			Sender:   1,
			Receiver: 2,
			Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
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
