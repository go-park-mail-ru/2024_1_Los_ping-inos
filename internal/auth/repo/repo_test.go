package repo

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	models "main.go/internal/auth"
	"main.go/internal/types"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"UserId", "Url", "CellNumber"})

	expect := []models.Image{
		{UserId: 1, Url: "url1", CellNumber: "0"},
		{UserId: 0, CellNumber: "1"},
		{UserId: 0, CellNumber: "2"},
		{UserId: 0, CellNumber: "3"},
		{UserId: 0, CellNumber: "4"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.UserId, item.Url, item.CellNumber)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT person_id, image_url, cell_number FROM person_image WHERE person_id = $1")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &ImageStorage{
		dbReader: db,
	}

	images, err := repo.Get(context.TODO(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(images, expect) {
		t.Errorf("expected %v, got %v", expect, images)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT person_id, image_url, cell_number FROM person_image WHERE person_id = $1")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	images, err = repo.Get(context.TODO(), 1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if len(images) != 0 {
		t.Errorf("expected empty array, got %v", images)
		return
	}
}

func TestCreatePersonInterest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"chto", "za", "govno"})

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO person_interest (person_id,interest_id) VALUES ($1,$2)")).
		WithArgs(1, 1).WillReturnRows(rows)

	repo := &InterestStorage{
		dbReader: db,
	}

	err = repo.CreatePersonInterests(context.TODO(), 1, []types.InterestID{1})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO person_interest (person_id,interest_id) VALUES ($1,$2)")).
		WithArgs(1, 1).WillReturnError(fmt.Errorf("repo error"))

	err = repo.CreatePersonInterests(context.TODO(), 1, []types.InterestID{1})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected no error, got %v", err)
	}

}
