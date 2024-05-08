package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"main.go/internal/auth"
	models "main.go/internal/auth"
	mocks "main.go/internal/auth/mocks"
	image "main.go/internal/image/protos/gen"
	"main.go/internal/types"
)

func TestGetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := mocks.NewMockPersonStorage(ctrl)

	PersonGetFilter := &auth.PersonGetFilter{
		SessionID: []string{"47300672-793e-4fd"},
		ID:        []types.UserID{1},
		Name:      "nikola_kwas",
	}

	persons := []*auth.Person{
		{
			ID:        types.UserID(1),
			Name:      "nikola_kwas",
			Email:     "nikola_kwas",
			SessionID: "47300672-793e-4fd",
		},
	}

	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter).Return(persons, nil)

	mockImage := mocks.NewMockImageClient(ctrl)

	imageResponce := &image.GetImageResponce{
		Url: "http://localhost",
	}

	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "0"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "1"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "2"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "3"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "4"}).Return(imageResponce, nil)

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

	core := UseCase{personStorage: mockObj, grpcClient: mockImage, interestStorage: mockInterest}

	testTable := []auth.ProfileGetParams{
		{
			ID:        []types.UserID{1},
			SessionID: []string{"47300672-793e-4fd"},
			Name:      "nikola_kwas",
		},
	}

	for _, curr := range testTable {
		prof, err := core.GetProfile(curr, context.TODO())
		if err != nil {
			t.Errorf("unexpected err result")
			return
		}
		if len(prof) == 0 {
			t.Errorf("unexpected result")
		}
		if prof[0].Name != curr.Name {
			t.Errorf("unexpected result")
		}
	}
}

func TestUpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := mocks.NewMockPersonStorage(ctrl)

	PersonGetFilter := &auth.PersonGetFilter{
		ID: []types.UserID{1},
	}

	hashedPassword, _ := hashPassword("password")

	persons := []*auth.Person{
		{
			ID:          types.UserID(1),
			Name:        "nikola_kwas",
			Email:       "nikola_kwas",
			SessionID:   "47300672-793e-4fd",
			Password:    hashedPassword,
			Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			Description: "hehehehaw",
		},
	}

	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter).Return(persons, nil)
	mockObj.EXPECT().Update(gomock.Any(), persons[0]).Return(nil)

	core := UseCase{personStorage: mockObj}

	testTable := []struct {
		UID  types.UserID
		prof auth.ProfileUpdateRequest
	}{
		{
			UID: types.UserID(1),
			prof: auth.ProfileUpdateRequest{
				Name:        "nikola_kwa",
				Email:       "nikola_kds",
				Password:    "password",
				Birthday:    "20011111",
				Description: "new description",
			},
		},
	}

	for _, curr := range testTable {
		err := core.UpdateProfile(curr.UID, curr.prof, context.TODO())
		if err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
	}
}

func TestDeleteProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := mocks.NewMockPersonStorage(ctrl)

	mockObj.EXPECT().Delete(gomock.Any(), types.UserID(1)).Return(nil)

	core := UseCase{personStorage: mockObj}

	testTable := []struct {
		UID types.UserID
	}{
		{
			UID: types.UserID(1),
		},
	}

	for _, curr := range testTable {
		//err := core.UpdateProfile(curr.UID, curr.prof, context.TODO())
		err := core.DeleteProfile(curr.UID, context.TODO())
		if err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
	}
}

func TestHandleInterests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	interestsBefore := []*models.Interest{
		{
			ID:   1,
			Name: "foo",
		},
		{
			ID:   2,
			Name: "bar",
		},
	}

	interestsAfter := []*models.Interest{
		{
			ID:   3,
			Name: "nice one",
		},
	}

	interestGetFilter := &models.InterestGetFilter{
		Name: []string{"nice one"},
	}

	mockObj := mocks.NewMockInterestStorage(ctrl)
	mockObj.EXPECT().GetPersonInterests(gomock.Any(), types.UserID(1)).Return(interestsBefore, nil)
	mockObj.EXPECT().Get(gomock.All(), interestGetFilter).Return(interestsAfter, nil)
	mockObj.EXPECT().CreatePersonInterests(gomock.Any(), types.UserID(1), []types.InterestID{3}).Return(nil)
	mockObj.EXPECT().DeletePersonInterests(gomock.Any(), types.UserID(1), []types.InterestID{1, 2}).Return(nil)

	core := UseCase{interestStorage: mockObj}

	testTable := []struct {
		interests []string
		userID    types.UserID
	}{
		{
			interests: []string{"nice one"},
			userID:    types.UserID(1),
		},
	}

	for _, curr := range testTable {
		//err := core.DeleteProfile(curr.UID, context.TODO())
		err := core.handleInterests(curr.interests, curr.userID, context.TODO())
		if err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
	}
}

func TestGetMatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := mocks.NewMockPersonStorage(ctrl)

	//mockObj.EXPECT().Delete(gomock.Any(), types.UserID(1)).Return(nil)
	mockObj.EXPECT().GetMatch(gomock.Any(), types.UserID(1)).Return([]types.UserID{2}, nil)

	PersonGetFilter := &auth.PersonGetFilter{
		//SessionID: []string{"47300672-793e-4fd"},
		ID:   []types.UserID{2},
		Name: "nikola_kwas",
	}

	persons := []*auth.Person{
		{
			ID:        types.UserID(1),
			Name:      "nikola_kwas",
			Email:     "nikola_kwas",
			SessionID: "47300672-793e-4fd",
		},
	}

	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter).Return(persons, nil)

	mockImage := mocks.NewMockImageClient(ctrl)

	imageResponce := &image.GetImageResponce{
		Url: "http://localhost",
	}

	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "0"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "1"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "2"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "3"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "4"}).Return(imageResponce, nil)

	mockInterest := mocks.NewMockInterestStorage(ctrl)

	interests := []*models.Interest{}

	mockInterest.EXPECT().GetPersonInterests(gomock.Any(), types.UserID(1)).Return(interests, nil)

	core := UseCase{personStorage: mockObj, grpcClient: mockImage, interestStorage: mockInterest}

	testTable := []struct {
		UID      types.UserID
		name     string
		profiles []auth.Profile
	}{
		{
			UID:  types.UserID(1),
			name: "nikola_kwas",
			profiles: []auth.Profile{
				{
					ID:   types.UserID(1),
					Name: "nikola_kwas",
					// Interests: []*models.Interest{
					// 	{
					// 		ID:   1,
					// 		Name: "foo",
					// 	},
					// 	{
					// 		ID:   2,
					// 		Name: "bar",
					// 	},
					// },
					Photos: []models.ImageToSend{
						{
							Cell: "0",
							Url:  "http://localhost",
						},
						{
							Cell: "1",
							Url:  "http://localhost",
						},
						{
							Cell: "2",
							Url:  "http://localhost",
						},
						{
							Cell: "3",
							Url:  "http://localhost",
						},
						{
							Cell: "4",
							Url:  "http://localhost",
						},
					},
				},
			},
		},
	}

	for _, curr := range testTable {
		//err := core.DeleteProfile(curr.UID, context.TODO())
		profiles, err := core.GetMatches(curr.UID, curr.name, context.TODO())
		if err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		if !reflect.DeepEqual(profiles, curr.profiles) {
			t.Errorf("unexpected profiles result")
			t.Error(profiles)
			t.Error(curr.profiles)
			return
		}
	}
}
