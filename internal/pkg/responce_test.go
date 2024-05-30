package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"main.go/internal/auth"
)

func TestSendResponse(t *testing.T) {
	mockBody := auth.Profile{ID: 1}

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	SendResponse(rr, req, http.StatusOK, mockBody)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, "https://jimder.ru", rr.Header().Get("Access-Control-Allow-Origin"))

	expectedResponse, err := easyjson.Marshal(mockBody)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedResponse), rr.Body.String())
}

func TestSendSimpleResponse(t *testing.T) {
	mockBody := "test message"

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	SendSimpleResponse(rr, req, http.StatusOK, mockBody)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, "https://jimder.ru", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, mockBody, rr.Body.String())
}
