package storage

import (
	"github.com/DATA-DOG/go-sqlmock"
	models "main.go/db"
	"main.go/internal/types"
	"reflect"
	"testing"
	"time"
)

func TestPersonGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var person1ID types.UserID = 1

	rows := sqlmock.
		NewRows([]string{"id", "name", "birthday", "description", "location", "photo", "email", "password", "created_at", "premium", "likes_left", "session_id", "gender"})
	expect := []*models.Person{{ID: person1ID, Name: "oleg"}}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name, time.Now(), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}

	mock.
		ExpectQuery("SELECT * FROM person WHERE (1=1)").
		WithArgs().
		WillReturnRows(rows)

	repo := &PersonStorage{
		dbReader: db,
	}
	likes, err := repo.Get(0, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(likes, expect) {
		t.Errorf("results not match, want %v, have %v", expect, likes)
		return
	}

}
