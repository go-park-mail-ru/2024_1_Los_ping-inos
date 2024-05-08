package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"main.go/internal/image/mocks"
)

func TestGetImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mocks.NewMockImgStorage(mockCtrl)

	mockObj.EXPECT().Get(gomock.Any(), int64(1), "1").Return("image", nil)

	core := UseCase{imageStorage: mockObj}

	testTable := []struct {
		id    int64
		name  string
		image string
	}{
		{
			id:    1,
			name:  "1",
			image: "image",
		},
	}

	for _, curr := range testTable {
		image, err := core.GetImage(curr.id, curr.name, context.TODO())
		require.Equal(t, image, curr.image)
		require.NoError(t, err)
	}
}

// func TestAddImage(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	mockObj := mocks.NewMockImgStorage(mockCtrl)

// 	userImage := models.Image{
// 		UserId:     1,
// 		Url:        "image.com",
// 		CellNumber: "1",
// 		FileName:   "image",
// 	}

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	pr, pw := io.Pipe()

// 	go func() {
// 		defer pw.Close()
// 		defer writer.Close()

// 		// Create a file part in the multipart form
// 		fileWriter, _ := writer.CreateFormFile("file", "test_image.jpg")

// 		// Copy the content of the buffer to the file part
// 		io.Copy(fileWriter, body)
// 	}()

// 	mockObj.EXPECT().Add(gomock.Any(), userImage, pr).Return(nil)

// 	core := UseCase{imageStorage: mockObj}

// 	testTable := []struct {
// 		userImage models.Image
// 		img       multipart.File
// 	}{
// 		{
// 			userImage: userImage,
// 			img:       pr,
// 		},
// 	}

// 	for _, curr := range testTable {
// 		image, err := core.GetImage(curr.id, curr.name, context.TODO())
// 		require.Equal(t, image, curr.image)
// 		require.NoError(t, err)
// 	}
// }

func TestDeleteImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mocks.NewMockImgStorage(mockCtrl)

	mockObj.EXPECT().Get(gomock.Any(), int64(1), "1").Return("image", nil)

	core := UseCase{imageStorage: mockObj}

	testTable := []struct {
		id    int64
		name  string
		image string
	}{
		{
			id:    1,
			name:  "1",
			image: "image",
		},
	}

	for _, curr := range testTable {
		image, err := core.GetImage(curr.id, curr.name, context.TODO())
		require.Equal(t, image, curr.image)
		require.NoError(t, err)
	}
}
