package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	mocks "main.go/internal/auth/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	models "main.go/internal/auth"
	image "main.go/internal/image/protos/gen"
	"main.go/internal/types"
)

func TestGetName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	getFilter := &models.PersonGetFilter{ID: []types.UserID{1}}

	expected := []*models.Person{
		{
			Name: "Nikita",
		},
	}

	mockObj := mocks.NewMockPersonStorage(ctrl)
	firstCall := mockObj.EXPECT().Get(gomock.Any(), getFilter).Return(expected, nil)
	mockObj.EXPECT().Get(gomock.Any(), getFilter).After(firstCall).Return(nil, fmt.Errorf("repo_error"))

	core := UseCase{personStorage: mockObj}

	fmt.Printf("%v", expected)

	result, err := core.GetName(1, context.TODO())
	println(result)
	newResult := []*models.Person{
		{
			Name: result,
		},
	}
	if err != nil {
		t.Errorf("unexpected error %s", err)
		return
	}
	if !reflect.DeepEqual(*expected[0], *newResult[0]) {
		t.Errorf("wanted %v, had %v", *expected[0], *newResult[0])
		return
	}

	result, err = core.GetName(1, context.TODO())
	if err == nil {
		t.Errorf("wanted error")
		return
	}
	if result != "" {
		t.Errorf("unexpected result")
		return
	}

}

func TestIsAuthenticated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []*models.Session{
		{
			UID: 1,
			SID: "47300672-793e-4fd4-b2a6-a1f12f16b83d",
		},
		{
			UID: 2,
			SID: "9d12e9e1-8d53-4d8a-811a-9aceaff0bd50",
		},
		{
			UID: 5,
			SID: "fa764ba8-008c-440a-a6db-1d5ee32484f9",
		},
	}

	mockObj := mocks.NewMockSessionStorage(ctrl)

	mockObj.EXPECT().GetBySID(gomock.Any(), "47300672-793e-4fd4-b2a6-a1f12f16b83d").Return(expected[0], nil)
	mockObj.EXPECT().GetBySID(gomock.Any(), "9d12e9e1-8d53-4d8a-811a-9aceaff0bd50").Return(expected[1], nil)
	mockObj.EXPECT().GetBySID(gomock.Any(), "fa764ba8-008c-440a-a6db-1d5ee32484f9").Return(nil, fmt.Errorf("no such user"))
	mockObj.EXPECT().GetBySID(gomock.Any(), "fa76sdf8-008c-440a-a6db-1d5ee32484f9").Return(nil, fmt.Errorf("repo_error"))
	mockObj.EXPECT().GetBySID(gomock.Any(), "").Return(nil, fmt.Errorf("repo_error"))

	core := UseCase{sessionStorage: mockObj}

	testTable := []struct {
		UID    types.UserID
		SID    string
		hasErr bool
		result bool
	}{
		{
			UID:    1,
			SID:    "47300672-793e-4fd4-b2a6-a1f12f16b83d",
			hasErr: false,
			result: true,
		},
		{
			UID:    2,
			SID:    "9d12e9e1-8d53-4d8a-811a-9aceaff0bd50",
			hasErr: false,
			result: true,
		},
		{
			UID:    5,
			SID:    "fa764ba8-008c-440a-a6db-1d5ee32484f9",
			hasErr: false,
			result: false,
		},
		{
			UID:    -1,
			SID:    "fa76sdf8-008c-440a-a6db-1d5ee32484f9",
			hasErr: true,
			result: false,
		},
		{
			UID:    -1,
			SID:    "",
			hasErr: true,
			result: false,
		},
	}

	for _, curr := range testTable {
		_, result, err := core.IsAuthenticated(curr.SID, context.TODO())
		if curr.hasErr && err == nil {
			t.Errorf("unexpected err result")
			return
		}
		if result != curr.result {
			t.Errorf("unexpected result")
		}
	}
}

func TestLogin(t *testing.T) {
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

	hashedPassword, _ := hashPassword("qwertyqwerty")

	expected := []*models.Person{
		{
			ID:       1,
			Name:     "Sanya",
			Email:    "somemail@gmial.com",
			Password: hashedPassword,
		},
	}

	getFilter := &models.PersonGetFilter{Email: []string{"somemail@gmial.com"}}

	mockSQL := mocks.NewMockPersonStorage(ctrl)
	mockSQL.EXPECT().Get(gomock.Any(), getFilter).Return(expected, nil)

	mockREDIS := mocks.NewMockSessionStorage(ctrl)
	mockREDIS.EXPECT().CreateSession(gomock.Any(), types.UserID(1)).Return("predefined_session_id", nil)

	mockInterest := mocks.NewMockInterestStorage(ctrl)

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

	core := UseCase{sessionStorage: mockREDIS, personStorage: mockSQL, interestStorage: mockInterest, grpcClient: mockObj}

	testTable := []struct {
		ID       types.UserID
		Name     string
		Email    string
		Password string
		hasErr   bool
	}{
		{
			ID:       1,
			Name:     "Sanya",
			Email:    "somemail@gmial.com",
			Password: "qwertyqwerty",
			hasErr:   false,
		},
	}

	for _, curr := range testTable {
		hashedPassword, _ := hashPassword(curr.Password)
		t.Log(hashedPassword)
		t.Log(curr.Password)
		profile, sessionID, err := core.Login(curr.Email, curr.Password, context.TODO())
		if err != nil {
			t.Error("Failed to login", err)
			return
		}
		if curr.hasErr && err == nil {
			t.Errorf("unexpected err result")
			return
		}
		t.Logf(curr.Name)
		t.Logf(profile.Email)
		if profile.Name != curr.Name {
			t.Errorf("unexpected profile")
		}
		if sessionID == "" {
			t.Errorf("unexpected sessionID")
		}
	}
}

func TestRedgistration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hashedPassword, _ := hashPassword("qwertyqwerty")

	expected := []*models.Person{
		{
			ID:       1,
			Name:     "Sanya",
			Birthday: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:   "male",
			Email:    "somemail@gmial.com",
			Password: hashedPassword,
		},
	}

	getFilter := &models.PersonGetFilter{Email: []string{"somemail@gmial.com"}}

	mockSQL := mocks.NewMockPersonStorage(ctrl)

	mockSQL.EXPECT().AddAccount(gomock.Any(), expected[0].Name, "20010101",
		expected[0].Gender, expected[0].Email, "qwertyqwerty").Return(hashedPassword, nil)

	mockSQL.EXPECT().Get(gomock.Any(), getFilter).Return(expected, nil)

	mockInterest := mocks.NewMockInterestStorage(ctrl)

	interestFilter := &models.InterestGetFilter{Name: []string{"foo", "bar"}}

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

	mockInterest.EXPECT().Get(gomock.Any(), interestFilter).Return(interests, nil)
	mockInterest.EXPECT().CreatePersonInterests(gomock.Any(), expected[0].ID, []types.InterestID{1, 2}).Return(nil)

	mockREDIS := mocks.NewMockSessionStorage(ctrl)
	mockREDIS.EXPECT().CreateSession(gomock.Any(), types.UserID(1)).Return("predefined_session_id", nil)

	interest := &mocks.MockInterestStorage{}

	core := UseCase{sessionStorage: mockREDIS, personStorage: mockSQL, interestStorage: interest}

	testTable := []models.RegitstrationBody{
		{
			Name:      "Sanya",
			Birthday:  "20010101",
			Gender:    "male",
			Email:     "somemail@gmial.com",
			Password:  "qwertyqwerty",
			Interests: []string{"foo", "bar"},
		},
	}

	for _, curr := range testTable {
		//hashedPassword, _ := hashPassword(curr.Password)
		//t.Log(hashedPassword)
		//t.Log(curr.Password)
		profile, sessionID, err := core.Registration(curr, context.TODO())
		require.NoError(t, err)
		require.Equal(t, curr.Name, profile.Name)
		require.Equal(t, curr.Email, profile.Email)
		require.NotEmpty(t, sessionID)
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

	mockInterest := mocks.NewMockInterestStorage(ctrl)

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

	core := UseCase{grpcClient: mockObj, interestStorage: mockInterest}

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

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockREDIS := mocks.NewMockSessionStorage(ctrl)
	mockREDIS.EXPECT().DeleteSession(gomock.Any(), "predefined_session_id").Return(nil)
	mockREDIS.EXPECT().DeleteSession(gomock.Any(), "wrong_session_id").Return(fmt.Errorf("repo_error"))

	core := UseCase{sessionStorage: mockREDIS}

	testTable := []struct {
		SID    string
		hasErr bool
	}{
		{
			SID:    "predefined_session_id",
			hasErr: false,
		},
		{
			SID:    "wrong_session_id",
			hasErr: true,
		},
	}

	for _, curr := range testTable {
		//profile, sessionID, err := core.Login(curr.Email, curr.Password, context.TODO())
		err := core.Logout(curr.SID, context.TODO())
		if curr.hasErr && err == nil {
			t.Errorf("unexpected err result")
			return
		}
	}
}
