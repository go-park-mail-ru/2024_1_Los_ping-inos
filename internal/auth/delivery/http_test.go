package delivery

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailru/easyjson"
	"main.go/internal/auth"
	mocks "main.go/internal/auth/mocks"
	. "main.go/internal/logs"
	requests "main.go/internal/pkg"
	"main.go/internal/types"

	_ "github.com/lib/pq"
)

func getResponse(w *httptest.ResponseRecorder) (*requests.Response, error) {
	var response requests.Response

	body, _ := io.ReadAll(w.Body)
	err := easyjson.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("cant unmarshal jsone")
	}

	return &response, nil
}

func TestReadProfile(t *testing.T) {

	testCases := map[string]struct {
		method string
		params map[string]string
		result requests.Response
	}{
		"Bad method": {
			method: http.MethodPost,
			params: map[string]string{},
			result: requests.Response{Status: http.StatusMethodNotAllowed, Body: nil},
		},
		"No Film": {
			method: http.MethodGet,
			params: map[string]string{},
			result: requests.Response{Status: http.StatusBadRequest, Body: nil},
		},
		"Core error": {
			method: http.MethodGet,
			params: map[string]string{"id": "1"},
			result: requests.Response{Status: http.StatusInternalServerError, Body: nil},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCore := mocks.NewMockIUseCase(ctrl)
	mockCore.EXPECT().GetProfile(auth.ProfileGetParams{ID: []types.UserID{types.UserID(0)}, NeedEmail: true}, gomock.Any()).Return(nil, nil)

	logger := InitLog()
	logger.RequestID = int64(1)
	//contexted := context.WithValue(context.Background(), Logg, logger)

	api := AuthHandler{UseCase: mockCore}

	for _, curr := range testCases {
		r := httptest.NewRequest(curr.method, "/api/v1/profile", nil)
		q := r.URL.Query()
		for key, value := range curr.params {
			q.Add(key, value)
		}
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()

		api.ReadProfile(w, r)
		response, err := getResponse(w)
		if err != nil {
			t.Error(err)
			return
		}
		if response.Status != curr.result.Status {
			t.Errorf("expected status %d, got %d", curr.result.Status, response.Status)
		}
	}
}

// func TestUpdateProfile(t *testing.T) {

// 	testCases := map[string]struct {
// 		method string
// 		params map[string]string
// 		result requests.Response
// 	}{
// 		"Bad method": {
// 			method: http.MethodPost,
// 			params: map[string]string{},
// 			result: requests.Response{Status: http.StatusMethodNotAllowed, Body: nil},
// 		},
// 		"No Film": {
// 			method: http.MethodGet,
// 			params: map[string]string{},
// 			result: requests.Response{Status: http.StatusBadRequest, Body: nil},
// 		},
// 		"Core error": {
// 			method: http.MethodGet,
// 			params: map[string]string{"id": "1"},
// 			result: requests.Response{Status: http.StatusInternalServerError, Body: nil},
// 		},
// 	}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockCore := mocks.NewMockIUseCase(ctrl)
// 	mockCore.EXPECT().GetProfile(auth.ProfileGetParams{ID: []types.UserID{types.UserID(0)}, NeedEmail: true}, gomock.Any()).Return(nil, nil)

// 	logger := InitLog()
// 	logger.RequestID = int64(1)
// 	//contexted := context.WithValue(context.Background(), Logg, logger)

// 	api := AuthHandler{UseCase: mockCore}

// 	for _, curr := range testCases {
// 		r := httptest.NewRequest(curr.method, "/api/v1/profile", nil)
// 		q := r.URL.Query()
// 		for key, value := range curr.params {
// 			q.Add(key, value)
// 		}
// 		r.URL.RawQuery = q.Encode()
// 		w := httptest.NewRecorder()

// 		api.UpdateProfile(w, r)
// 		response, err := getResponse(w)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		if response.Status != curr.result.Status {
// 			t.Errorf("expected status %d, got %d", curr.result.Status, response.Status)
// 		}
// 	}
// }

// func TestRegistrationHandler(t *testing.T) {

// 	testCases := map[string]struct {
// 		method string
// 		params map[string]string
// 		result requests.Response
// 	}{
// 		"Bad method": {
// 			method: http.MethodPost,
// 			params: map[string]string{},
// 			result: requests.Response{Status: http.StatusMethodNotAllowed, Body: nil},
// 		},
// 		"No Film": {
// 			method: http.MethodGet,
// 			params: map[string]string{},
// 			result: requests.Response{Status: http.StatusBadRequest, Body: nil},
// 		},
// 		"Core error": {
// 			method: http.MethodGet,
// 			params: map[string]string{"id": "1"},
// 			result: requests.Response{Status: http.StatusInternalServerError, Body: nil},
// 		},
// 	}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockCore := mocks.NewMockIUseCase(ctrl)
// 	mockCore.EXPECT().GetProfile(auth.ProfileGetParams{ID: []types.UserID{types.UserID(0)}, NeedEmail: true}, gomock.Any()).Return(nil, nil)

// 	logger := InitLog()
// 	logger.RequestID = int64(1)
// 	//contexted := context.WithValue(context.Background(), Logg, logger)

// 	api := AuthHandler{UseCase: mockCore}

// 	for _, curr := range testCases {
// 		r := httptest.NewRequest(curr.method, "/api/v1/profile", nil)
// 		q := r.URL.Query()
// 		for key, value := range curr.params {
// 			q.Add(key, value)
// 		}
// 		r.URL.RawQuery = q.Encode()
// 		w := httptest.NewRecorder()

// 		api.RegistrationHandler()
// 		response, err := getResponse(w)
// 		if err != nil {
// 			t.Error(err)
// 			return
// 		}
// 		if response.Status != curr.result.Status {
// 			t.Errorf("expected status %d, got %d", curr.result.Status, response.Status)
// 		}
// 	}
// }
