package storage

import (
	"github.com/DATA-DOG/go-sqlmock"
	models "main.go/db"
	"main.go/internal/types"
	"reflect"
	"testing"
)

func TestLikeGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var person1ID types.UserID = 1

	rows := sqlmock.
		NewRows([]string{"person_one_id", "person_two_id"})
	expect := []*models.Like{{person1ID, person1ID + 1}}
	for _, item := range expect {
		rows = rows.AddRow(item.Person1, item.Person2)
	}

	mock.
		ExpectQuery("SELECT * FROM \"like\" WHERE (person_one_id = $1)").
		WithArgs(person1ID).
		WillReturnRows(rows)

	repo := &LikeStorage{
		dbReader: db,
	}
	likes, err := repo.Get(0, &models.LikeGetFilter{Person1: &person1ID})

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
