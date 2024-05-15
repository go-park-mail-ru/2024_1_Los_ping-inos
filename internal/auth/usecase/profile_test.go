package usecase

import (
	"context"
	"fmt"
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

	PersonGetFilter := []*auth.PersonGetFilter{
		{
			SessionID: []string{"47300672-793e-4fd"},
			ID:        []types.UserID{1},
			Name:      "nikola_kwas",
		},
		{
			SessionID: []string{"47300672-793e-sdf"},
			ID:        []types.UserID{2},
			Name:      "nikola_kwas",
		},
		{
			SessionID: []string{"47300672-793e-kal"},
			ID:        []types.UserID{3},
			Name:      "nikola_kwas",
		},
	}

	persons := []*auth.Person{
		{
			ID:        types.UserID(1),
			Name:      "nikola_kwas",
			Email:     "nikola_kwas",
			SessionID: "47300672-793e-4fd",
		},
	}
	wrongPersons := []*auth.Person{
		{
			ID:        types.UserID(3),
			Name:      "nikola_kwas",
			Email:     "nikola_kwas",
			SessionID: "47300672-793e-kal",
		},
	}

	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[0]).Return(persons, nil)
	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[1]).Return(nil, fmt.Errorf("repo error"))
	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[2]).Return(wrongPersons, nil)

	mockImage := mocks.NewMockImageClient(ctrl)

	imageResponce := &image.GetImageResponce{
		Url: "http://localhost",
	}

	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "0"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "1"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "2"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "3"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(1), Cell: "4"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "0"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "1"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "2"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "3"}).Return(imageResponce, nil)
	mockImage.EXPECT().GetImage(gomock.Any(), &image.GetImageRequest{Id: int64(3), Cell: "4"}).Return(imageResponce, nil)

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
	mockInterest.EXPECT().GetPersonInterests(gomock.Any(), types.UserID(3)).Return(nil, fmt.Errorf("repo error"))

	core := UseCase{personStorage: mockObj, grpcClient: mockImage, interestStorage: mockInterest}

	testTable := []struct {
		profile auth.ProfileGetParams
		hasErr  bool
	}{
		{
			profile: auth.ProfileGetParams{

				SessionID: []string{"47300672-793e-4fd"},
				ID:        []types.UserID{1},
				Name:      "nikola_kwas",
			},
			hasErr: false,
		},
		{
			profile: auth.ProfileGetParams{

				SessionID: []string{"47300672-793e-sdf"},
				ID:        []types.UserID{2},
				Name:      "nikola_kwas",
			},
			hasErr: true,
		},
		{
			profile: auth.ProfileGetParams{

				SessionID: []string{"47300672-793e-kal"},
				ID:        []types.UserID{3},
				Name:      "nikola_kwas",
			},
			hasErr: true,
		},
	}

	for _, curr := range testTable {
		prof, err := core.GetProfile(curr.profile, context.TODO())
		if curr.hasErr && err == nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		if !curr.hasErr && err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		if !curr.hasErr && len(prof) == 0 {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		if !curr.hasErr && prof[0].Name != curr.profile.Name {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
	}
}

func TestUpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := mocks.NewMockPersonStorage(ctrl)

	PersonGetFilter := []*auth.PersonGetFilter{
		{
			ID: []types.UserID{1},
		},
		{
			ID: []types.UserID{2},
		},
		{
			ID: []types.UserID{3},
		},
		{
			ID: []types.UserID{4},
		},
		{
			ID: []types.UserID{5},
		},
		{
			ID: []types.UserID{6},
		},
		{
			ID: []types.UserID{7},
		},
	}

	hashedPassword, _ := hashPassword("password")

	persons := [][]*auth.Person{
		{
			{
				ID:          1,
				Name:        "nikola_kwas",
				Email:       "nikola_kwas",
				Password:    hashedPassword,
				Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				Description: "hehehehaw",
			},
		},
		{
			{
				ID:          2,
				Name:        "petro",
				Email:       "perto@sobaka.dom",
				Password:    hashedPassword,
				Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				Description: "hehehehaw",
			},
		},
		{
			{
				ID:          3,
				Name:        "petro",
				Email:       "perto@sobaka.dom",
				Password:    hashedPassword,
				Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				Description: "hehehehaw",
			},
		},
		{
			{
				ID:          4,
				Name:        "petro",
				Email:       "perto@sobaka.dom",
				Password:    hashedPassword,
				Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				Description: "hehehehaw",
			},
		},
		{
			{
				ID:          5,
				Name:        "petro",
				Email:       "perto@sobaka.dom",
				Password:    hashedPassword,
				Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				Description: "hehehehaw",
			},
		},
		{
			{
				ID:          6,
				Name:        "petro",
				Email:       "perto@sobaka.dom",
				Password:    hashedPassword,
				Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				Description: "hehehehaw",
			},
		},
		{
			{
				ID:          7,
				Name:        "petro",
				Email:       "perto@sobaka.dom",
				Password:    hashedPassword,
				Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
				Description: "hehehehaw",
			},
		},
	}
	updates := []*auth.Person{
		{
			ID:          1,
			Name:        "nikola_kwa",
			Email:       "nikola_kds",
			Password:    hashedPassword,
			Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			Description: "new description",
		},
		{
			ID:          7,
			Name:        "nikola_kwa",
			Email:       "nikola_kds",
			Password:    hashedPassword,
			Birthday:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
			Description: "new description",
		},
	}

	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[0]).Return(persons[0], nil)
	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[1]).Return(nil, fmt.Errorf("repo error"))
	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[2]).Return(persons[2], nil)
	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[3]).Return(persons[3], nil)
	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[4]).Return(persons[4], nil)
	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[5]).Return(persons[5], nil)
	mockObj.EXPECT().Get(gomock.Any(), PersonGetFilter[6]).Return(persons[6], nil)
	mockObj.EXPECT().Update(gomock.Any(), *updates[0]).Return(nil)
	mockObj.EXPECT().Update(gomock.Any(), *updates[1]).Return(fmt.Errorf("repo error"))

	core := UseCase{personStorage: mockObj}

	testTable := []struct {
		UID    types.UserID
		prof   *auth.ProfileUpdateRequest
		hasErr bool
	}{
		{
			UID: types.UserID(1),
			prof: &auth.ProfileUpdateRequest{
				Name:        "nikola_kwa",
				Email:       "nikola_kds",
				OldPassword: "password",
				Birthday:    "01.01.2001",
				Description: "new description",
			},
			hasErr: false,
		},
		{
			UID: types.UserID(2),
			prof: &auth.ProfileUpdateRequest{
				Name:        "nikola_kwa",
				Email:       "nikola_kds",
				OldPassword: "password",
				Birthday:    "01.01.2001",
				Description: "new description",
			},
			hasErr: true,
		},
		{
			UID: types.UserID(3),
			prof: &auth.ProfileUpdateRequest{
				Name:        "nikola_kwa",
				Email:       "nikola_kds",
				OldPassword: "passwordd",
				Birthday:    "01.01.2001",
				Description: "new description",
			},
			hasErr: true,
		},
		{
			UID: types.UserID(4),
			prof: &auth.ProfileUpdateRequest{
				Name:        "nikola_kwa",
				Email:       "nikola_kds",
				OldPassword: "password",
				Birthday:    "20010101",
				Description: "new description",
			},
			hasErr: true,
		},
		{
			UID: types.UserID(5),
			prof: &auth.ProfileUpdateRequest{
				Name:        "nikola_kwa",
				Email:       "nikola_kds",
				OldPassword: "password",
				Birthday:    "20010101",
				Description: "new description",
				Interests:   []string{"footba", "goo"},
			},
			hasErr: true,
		},
		{
			UID: types.UserID(6),
			prof: &auth.ProfileUpdateRequest{
				Name:        "nikola_kwa",
				Password:    "newpassword",
				OldPassword: "passworddsdf",
				Birthday:    "01.01.2001",
				Description: "new description",
			},
			hasErr: true,
		},
		{
			UID: types.UserID(7),
			prof: &auth.ProfileUpdateRequest{
				Name:        "nikola_kwa",
				Email:       "nikola_kds",
				OldPassword: "password",
				Birthday:    "01.01.2001",
				Description: "new description",
			},
			hasErr: true,
		},
	}

	for _, curr := range testTable {
		err := core.UpdateProfile(curr.UID, *curr.prof, context.TODO())
		if curr.hasErr && err == nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		// if err != nil {
		// 	t.Errorf("unexpected err result")
		// 	t.Error(err)
		// 	return
		// }
		if !curr.hasErr && err != nil {
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

	mockObj.EXPECT().GetMatch(gomock.Any(), types.UserID(1)).Return([]types.UserID{2}, nil)
	mockObj.EXPECT().GetMatch(gomock.Any(), types.UserID(2)).Return(nil, fmt.Errorf("repo error"))
	mockObj.EXPECT().GetMatch(gomock.Any(), types.UserID(3)).Return([]types.UserID{}, nil)

	PersonGetFilter := &auth.PersonGetFilter{
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
		hasErr   bool
	}{
		{
			UID:  types.UserID(1),
			name: "nikola_kwas",
			profiles: []auth.Profile{
				{
					ID:   types.UserID(1),
					Name: "nikola_kwas",
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
					Interests: []*models.Interest{},
				},
			},
			hasErr: false,
		},
		{
			UID:    types.UserID(2),
			hasErr: true,
		},
		{
			UID:      types.UserID(3),
			hasErr:   false,
			profiles: []auth.Profile{},
		},
	}

	for _, curr := range testTable {
		profiles, err := core.GetMatches(curr.UID, curr.name, context.TODO())
		if curr.hasErr && err == nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		if !curr.hasErr && err != nil {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		if !curr.hasErr && !reflect.DeepEqual(profiles, curr.profiles) {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
		//require.Equal(t, profiles, curr.profiles)
	}
}
