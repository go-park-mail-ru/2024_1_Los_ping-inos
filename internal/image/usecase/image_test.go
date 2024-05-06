package usecase

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
}
