package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	models "main.go/internal/feed"
	mocks "main.go/internal/feed/mocks"
	image "main.go/internal/image/protos/gen"
	"main.go/internal/types"
)

func TestNewUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockPostgresStorage(ctrl)
	mockWs := mocks.NewMockWebSocStorage(ctrl)
	mockGrpc := mocks.NewMockImageClient(ctrl)

	useCase := New(mockStorage, mockWs, mockGrpc)

	if useCase.storage == nil {
		t.Error("personStorage should not be nil")
	}
	if useCase.ws == nil {
		t.Error("sessionStorage should not be nil")
	}
	if useCase.grpcClient == nil {
		t.Error("interestStorage should not be nil")
	}
}

func TestGetUserCards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := mocks.NewMockImageClient(ctrl)

	imageResponce := &image.GetImageResponce{
		Url: "http://localhost",
	}

	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "0"}).Return(imageResponce, nil)
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "1"}).Return(imageResponce, nil)
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "2"}).Return(imageResponce, nil)
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "3"}).Return(imageResponce, nil)
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "4"}).Return(imageResponce, nil)

	mockInterest := mocks.NewMockPostgresStorage(ctrl)

	interests := []*models.Interest{
		{
			ID:   1,
			Name: "foo",
		},
		{
			ID:   2,
			Name: "bar",
		},
	}

	mockInterest.EXPECT().GetPersonInterests(gomock.Any(), types.UserID(1)).Return(interests, nil)

	core := UseCase{grpcClient: mockObj, storage: mockInterest}

	testTable := []*models.Person{
		{
			ID:       1,
			Name:     "Sanya",
			Birthday: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:   "male",
			Email:    "somemail@gmial.com",
		},
	}

	newInterests := [][]*models.Interest{
		{
			{
				ID:   1,
				Name: "foo",
			},
			{
				ID:   2,
				Name: "bar",
			},
		},
	}

	newImage := [][]models.Image{
		{
			{
				UserId:     1,
				Url:        "http://localhost",
				CellNumber: "0",
			},
			{
				UserId:     1,
				Url:        "http://localhost",
				CellNumber: "1",
			},
			{
				UserId:     1,
				Url:        "http://localhost",
				CellNumber: "2",
			},
			{
				UserId:     1,
				Url:        "http://localhost",
				CellNumber: "3",
			},
			{
				UserId:     1,
				Url:        "http://localhost",
				CellNumber: "4",
			},
		},
	}

	interes, imm, err := core.getUserCards(testTable, context.TODO())
	require.NoError(t, err)
	require.Equal(t, newInterests, interes)
	require.Equal(t, newImage, imm)
}

func TestGetCards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := mocks.NewMockImageClient(ctrl)

	imageResponce := &image.GetImageResponce{
		Url: "http://localhost",
	}

	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "0"}).Return(imageResponce, nil)
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "1"}).Return(imageResponce, nil)
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "2"}).Return(imageResponce, nil)
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "3"}).Return(imageResponce, nil)
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "4"}).Return(imageResponce, nil)
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "0"}).Return(nil, fmt.Errorf("repo error"))
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "1"}).Return(nil, fmt.Errorf("repo error"))
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "2"}).Return(nil, fmt.Errorf("repo error"))
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "3"}).Return(nil, fmt.Errorf("repo error"))
	mockObj.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "4"}).Return(nil, fmt.Errorf("repo error"))

	mockSQL := mocks.NewMockPostgresStorage(ctrl)

	interests := []*models.Interest{
		{
			ID:   1,
			Name: "foo",
		},
		{
			ID:   2,
			Name: "bar",
		},
	}

	mockSQL.EXPECT().GetPersonInterests(gomock.Any(), types.UserID(1)).Return(interests, nil)
	mockSQL.EXPECT().GetPersonInterests(gomock.Any(), types.UserID(3)).Return(nil, fmt.Errorf("repo error"))

	persons := []*models.Person{
		{
			ID:       1,
			Name:     "Sanya",
			Birthday: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:   "male",
			Email:    "somemail@gmial.com",
		},
	}
	wrongPersons := []*models.Person{
		{
			ID:       3,
			Name:     "Sanya",
			Birthday: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:   "male",
			Email:    "somemail@gmial.com",
		},
	}

	mockSQL.EXPECT().GetFeed(gomock.Any(), types.UserID(1)).Return(persons, nil)
	mockSQL.EXPECT().GetFeed(gomock.Any(), types.UserID(2)).Return(nil, fmt.Errorf("repo error"))
	mockSQL.EXPECT().GetFeed(gomock.Any(), types.UserID(3)).Return(wrongPersons, nil)

	core := UseCase{grpcClient: mockObj, storage: mockSQL}

	newInterests := [][]*models.Interest{
		{
			{
				ID:   1,
				Name: "foo",
			},
			{
				ID:   2,
				Name: "bar",
			},
		},
	}

	newImage := [][]models.ImageToSend{
		{
			{
				//UserId:     1,
				Url:  "http://localhost",
				Cell: "0",
			},
			{
				Url:  "http://localhost",
				Cell: "1",
			},
			{

				Url:  "http://localhost",
				Cell: "2",
			},
			{

				Url:  "http://localhost",
				Cell: "3",
			},
			{

				Url:  "http://localhost",
				Cell: "4",
			},
		},
	}

	testTable := []struct {
		cards  []models.Card
		hasErr bool
	}{
		{
			cards: []models.Card{
				{
					ID:        1,
					Name:      "Sanya",
					Birthday:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					Interests: newInterests[0],
					Photos:    newImage[0],
				},
			},
			hasErr: false,
		},
		{
			cards: []models.Card{
				{
					ID:        2,
					Name:      "Sanya",
					Birthday:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					Interests: newInterests[0],
					Photos:    newImage[0],
				},
			},
			hasErr: true,
		},
		{
			cards: []models.Card{
				{
					ID:        3,
					Name:      "Sanya",
					Birthday:  time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					Interests: newInterests[0],
					Photos:    newImage[0],
				},
			},
			hasErr: true,
		},
	}

	//interes, imm, err := core.getUserCards(testTable, context.TODO())
	for _, curr := range testTable {
		result, err := core.GetCards(curr.cards[0].ID, context.TODO())
		if curr.hasErr && err == nil {
			t.Errorf("unexpected err result")
			return
		}
		if !curr.hasErr && err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		if !curr.hasErr && !reflect.DeepEqual(result, curr.cards) {
			t.Errorf("unexpected err result")
			return
		}
	}
}

func TestCreateLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := mocks.NewMockPostgresStorage(ctrl)
	mockObj.EXPECT().CreateLike(gomock.Any(), types.UserID(1), types.UserID(2)).Return(nil)

	core := UseCase{storage: mockObj}

	testTable := []struct {
		profile1 types.UserID
		profile2 types.UserID
	}{
		{
			profile1: types.UserID(1),
			profile2: types.UserID(2),
		},
	}

	for _, curr := range testTable {
		err := core.CreateLike(curr.profile1, curr.profile2, context.TODO())
		if err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
	}
}

func TestGetChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chat := []models.Message{
		{
			Id:       1,
			Data:     "just got a new glock, its fif'een",
			Sender:   1,
			Receiver: 2,
			Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	mockObj := mocks.NewMockPostgresStorage(ctrl)
	mockObj.EXPECT().GetChat(gomock.Any(), types.UserID(1), types.UserID(2)).Return(chat, nil)

	core := UseCase{storage: mockObj}

	testTable := []struct {
		profile1 types.UserID
		profile2 types.UserID
		chat     []models.Message
	}{
		{
			profile1: types.UserID(1),
			profile2: types.UserID(2),
			chat: []models.Message{
				{
					Id:       1,
					Data:     "just got a new glock, its fif'een",
					Sender:   1,
					Receiver: 2,
					Time:     time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, curr := range testTable {
		chat, err := core.GetChat(context.TODO(), curr.profile1, curr.profile2)
		if err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		require.Equal(t, curr.chat, chat)
	}
}

func TestSaveMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chat := []models.MessageToReceive{
		{
			Id:       1,
			Data:     "just got a new glock, its fif'een",
			Sender:   1,
			Receiver: 2,
			Time:     20010101,
		},
	}

	mockObj := mocks.NewMockPostgresStorage(ctrl)
	mockObj.EXPECT().CreateMessage(gomock.Any(), chat[0]).Return(&chat[0], nil)

	core := UseCase{storage: mockObj}

	testTable := []struct {
		profile1 types.UserID
		profile2 types.UserID
		chat     []models.MessageToReceive
	}{
		{
			profile1: types.UserID(1),
			profile2: types.UserID(2),
			chat: []models.MessageToReceive{
				{
					Id:       1,
					Data:     "just got a new glock, its fif'een",
					Sender:   1,
					Receiver: 2,
					Time:     20010101,
				},
			},
		},
	}

	for _, curr := range testTable {
		chat, err := core.SaveMessage(context.TODO(), curr.chat[0])
		if err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		require.Equal(t, &curr.chat[0], chat)
	}
}

func TestCreateClaim(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	claim := []models.Claim{
		{
			TypeID:     1,
			SenderID:   1,
			ReceiverID: 2,
		},
	}

	mockObj := mocks.NewMockPostgresStorage(ctrl)
	mockObj.EXPECT().CreateClaim(gomock.Any(), claim[0]).Return(nil)

	core := UseCase{storage: mockObj}

	testTable := []struct {
		typeID     int64
		senderID   int64
		recieverID int64
	}{
		{
			typeID:     1,
			senderID:   1,
			recieverID: 2,
		},
	}

	for _, curr := range testTable {
		err := core.CreateClaim(context.TODO(), curr.typeID, curr.senderID, curr.recieverID)
		if err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
	}
}

func TestGetClaims(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	claim := []models.PureClaim{
		{
			Id:    1,
			Title: "whya",
		},
		{
			Id:    2,
			Title: "sdf",
		},
		{
			Id:    3,
			Title: "WHATTHEFUUUUUUUUUUUUUUUUUU",
		},
	}

	mockObj := mocks.NewMockPostgresStorage(ctrl)
	mockObj.EXPECT().GetAllClaims(gomock.Any()).Return(claim, nil)

	core := UseCase{storage: mockObj}

	result, err := core.GetClaims(context.TODO())
	if err != nil {
		t.Errorf("unexpected err result")
		t.Error(err)
		return
	}
	require.Equal(t, claim, result)
}
