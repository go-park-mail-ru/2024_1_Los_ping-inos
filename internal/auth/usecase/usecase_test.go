package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	mocks "main.go/internal/auth/mocks"

	"github.com/golang/mock/gomock"
	models "main.go/internal/auth"
	"main.go/internal/types"
)

// models "main.go/internal/auth"

// func TestGetAccount(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	testTable := []struct {
// 		name        string
// 		inputEmail  string
// 		inputPasswd string
// 		expSID      string
// 		expUID      int64
// 		expErr      error
// 	}{
// 		{
// 			name:        "successfully logs in existing user",
// 			inputEmail:  "johndoe@mail.com",
// 			inputPasswd: "passw0rd!",
// 			expSID:      "abcd1234efgh5678",
// 			expUID:      1,
// 			expErr:      nil,
// 		},
// 		{
// 			name:        "fails due to invalid credentials",
// 			inputEmail:  "nonexistinguser@mail.com",
// 			inputPasswd: "invalidpasswd",
// 			expSID:      "",
// 			expUID:      -1,
// 			expErr:      errors.New("incorrect username or password"),
// 		},
// 	}

// 	for _, tt := range testTable {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Set up mock PersonStorage
// 			ps := mocks.NewMockPersonStorage(ctrl)

// 			// Populate the mock data
// 			var user models.Person
// 			if tt.inputEmail == "johndoe@mail.com" {
// 				hashedPasswd, _ := hashPassword(tt.inputPasswd)
// 				user = models.Person{
// 					ID:          1,
// 					Name:        "John Doe",
// 					Birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
// 					Gender:      "Male",
// 					Email:       tt.inputEmail,
// 					Description: "",
// 					Password:    hashedPasswd,
// 				}
// 			}

// 			getFilter := &models.PersonGetFilter{Email: []string{tt.inputEmail}}
// 			ps.EXPECT().Get(gomock.Any(), getFilter).Return([]*models.Person{&user}, nil)

// 			ah := &UseCase{
// 				personStorage: ps,
// 			}

// 			user, gotErr := ah.GetProfile(tt.inputEmail, tt.inputPasswd, 1)

// 			assert.Equal(t, tt.expSID, gotSID)
// 			assert.Equal(t, tt.expUID, gotUID)
// 			assert.Equal(t, tt.expErr, gotErr)
// 		})
// 	}
// }

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
