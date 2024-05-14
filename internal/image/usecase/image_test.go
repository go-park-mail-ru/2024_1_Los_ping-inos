package usecase

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	models "main.go/internal/image"
	"main.go/internal/image/mocks"
)

func TestGetImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mocks.NewMockImgStorage(mockCtrl)

	mockObj.EXPECT().Get(gomock.Any(), int64(1), "1").Return("image", nil)
	mockObj.EXPECT().Get(gomock.Any(), int64(2), "1").Return("", fmt.Errorf("repo error"))
	mockObj.EXPECT().Get(gomock.Any(), int64(3), "1").Return("", nil)

	core := UseCase{imageStorage: mockObj}

	testTable := []struct {
		id     int64
		name   string
		image  string
		hasErr bool
	}{
		{
			id:     1,
			name:   "1",
			image:  "image",
			hasErr: false,
		},
		{
			id:     2,
			name:   "1",
			image:  "image",
			hasErr: true,
		},
		{
			id:     3,
			name:   "1",
			image:  "image",
			hasErr: true,
		},
	}

	for _, curr := range testTable {
		image, err := core.GetImage(curr.id, curr.name, context.TODO())
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
		if !curr.hasErr && image != curr.image {
			t.Errorf("unexpected err result")
			t.Error(err)
			return
		}
	}
}

func TestAddImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rs := os.File{}

	userImage := []models.Image{
		{
			UserId:     1,
			Url:        "image.com",
			CellNumber: "1",
			FileName:   "image",
		},
		{
			UserId:     2,
			Url:        "image.com",
			CellNumber: "1",
			FileName:   "image",
		},
	}

	mockObj := mocks.NewMockImgStorage(mockCtrl)

	mockObj.EXPECT().Add(gomock.Any(), userImage[0], &rs).Return(nil)
	mockObj.EXPECT().Add(gomock.Any(), userImage[1], &rs).Return(fmt.Errorf("repo error"))

	core := UseCase{imageStorage: mockObj}

	err := core.AddImage(userImage[0], &rs, context.TODO())
	if err != nil {
		t.Errorf("unexpected err result")
		t.Error(err)
		return
	}
	err = core.AddImage(userImage[1], &rs, context.TODO())
	if err == nil {
		t.Errorf("unexpected err result")
		t.Error(err)
		return
	}
}

func TestDeleteImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userImage := []models.Image{
		{
			UserId:     1,
			Url:        "image.com",
			CellNumber: "1",
			FileName:   "image",
		},
		{
			UserId:     2,
			Url:        "image.com",
			CellNumber: "1",
			FileName:   "image",
		},
	}

	mockObj := mocks.NewMockImgStorage(mockCtrl)

	mockObj.EXPECT().Delete(gomock.Any(), userImage[0]).Return(nil)
	mockObj.EXPECT().Delete(gomock.Any(), userImage[1]).Return(fmt.Errorf("repo error"))

	core := UseCase{imageStorage: mockObj}

	testTable := []struct {
		userImage models.Image
		hasErr    bool
	}{
		{
			userImage: userImage[0],
			hasErr:    false,
		},
		{
			userImage: userImage[1],
			hasErr:    true,
		},
	}

	for _, curr := range testTable {
		//image, err := core.GetImage(curr.id, curr.name, context.TODO())
		err := core.DeleteImage(curr.userImage, context.TODO())
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
	}
}

func TestNewImageUseCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := mocks.NewMockImgStorage(mockCtrl)

	core := UseCase{imageStorage: mockObj}

	testTable := []struct {
		istore models.ImgStorage
	}{
		{
			istore: mockObj,
		},
	}

	for _, curr := range testTable {
		UseCase := NewImageUseCase(curr.istore)
		require.Equal(t, &core, UseCase)
	}
}

func TestGetClient(t *testing.T) {
	port := ":50051"
	client, err := GetClient(port)
	if err != nil {
		t.Errorf("GetClient() error: %v", err)
	}
	if client == nil {
		t.Errorf("GetClient() returned nil client")
	}
}
