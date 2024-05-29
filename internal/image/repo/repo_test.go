package repo

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"main.go/internal/image"
	. "main.go/internal/logs"
)

//func TestGetImageRepo_Success(t *testing.T) {
//	db, err := sql.Open("postgres", "your_database_connection_string")
//	if err != nil {
//		t.Fatalf("Error opening database: %v", err)
//	}
//	defer db.Close()
//
//	storage, err := GetImageRepo("your_database_connection_string")
//	if err != nil {
//		t.Fatalf("Error creating ImageStorage: %v", err)
//	}
//
//	if storage.dbReader == nil {
//		t.Error("ImageStorage dbReader is nil")
//	}
//}

func TestNewImageStorage(t *testing.T) {
	db, err := sql.Open("postgres", "user=your_user password=your_password dbname=your_db sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	storage := NewImageStorage(db)
	if storage == nil {
		t.Fatalf("Failed to create ImageStorage instance")
	}
}

//func TestAdd(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("cant create mock: %s", err)
//	}
//	defer db.Close()
//
//	selectRow := "INSERT INTO person_image (person_id, image_url, cell_number) VALUES ($1, $2, $3) ON CONFLICT (person_id, cell_number) DO UPDATE SET image_url = EXCLUDED.image_url;"
//
//	logger := InitLog()
//	logger.RequestID = int64(1)
//
//	mock.ExpectExec(
//		regexp.QuoteMeta(selectRow)).
//		WithArgs(1, "1", "1").WillReturnResult(sqlmock.NewResult(0, 1))
//
//	repo := &ImageStorage{
//		dbReader: db,
//	}
//
//	contexted := context.WithValue(context.Background(), Logg, logger)
//
//	image := image.Image{
//		UserId:     1,
//		Url:        "1",
//		CellNumber: "1",
//		FileName:   "1",
//	}
//
//	//rs := os.File{}
//
//	err = repo.Add(contexted, image, nil)
//	if err != nil {
//		t.Errorf("repo error: %s", err)
//	}
//
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//
//	mock.ExpectExec(
//		regexp.QuoteMeta(selectRow)).
//		WithArgs(1, "1", "1").WillReturnError(fmt.Errorf("repo error"))
//
//	err = repo.Add(contexted, image, nil)
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//
//	if err == nil {
//		t.Errorf("repo error: %s", err)
//	}
//}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"UserId", "Url", "CellNumber"})

	expect := []image.Image{
		{UserId: 1, Url: "url1", CellNumber: "0"},
		{UserId: 0, CellNumber: "1"},
		{UserId: 0, CellNumber: "2"},
		{UserId: 0, CellNumber: "3"},
		{UserId: 0, CellNumber: "4"},
	}

	image := "url1"

	for _, item := range expect {
		rows = rows.AddRow(item.UserId, item.Url, item.CellNumber)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT person_id, image_url, cell_number FROM person_image WHERE person_id = $1 AND cell_number = $2")).
		WithArgs(1, "0").
		WillReturnRows(rows)

	print(rows)

	repo := &ImageStorage{
		dbReader: db,
	}

	images, err := repo.Get(context.TODO(), 1, "0")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(images, image) {
		t.Errorf("expected %v, got %v", image, images)
		print(images)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT person_id, image_url, cell_number FROM person_image WHERE person_id = $1 AND cell_number = $2")).
		WithArgs(1, "0").
		WillReturnError(fmt.Errorf("db_error"))

	images, err = repo.Get(context.TODO(), 1, "0")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if images != "" {
		t.Errorf("expected empty array, got %v", images)
		return
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	selectRow := "DELETE FROM person_image WHERE person_id = $1 AND cell_number = $2"

	logger := InitLog()
	logger.RequestID = int64(1)

	mock.ExpectExec(
		regexp.QuoteMeta(selectRow)).
		WithArgs(1, "1").WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &ImageStorage{
		dbReader: db,
	}

	contexted := context.WithValue(context.Background(), Logg, logger)

	image := image.Image{
		UserId:     1,
		Url:        "1",
		CellNumber: "1",
		FileName:   "1",
	}

	err = repo.Delete(contexted, image)
	if err != nil {
		t.Errorf("repo error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectExec(
		regexp.QuoteMeta(selectRow)).
		WithArgs(1, "1").WillReturnError(fmt.Errorf("repo err"))

	err = repo.Delete(contexted, image)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("repo error: %s", err)
	}
}
