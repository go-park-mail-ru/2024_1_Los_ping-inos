package repo

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	models "main.go/internal/feed"
	. "main.go/internal/logs"
	"main.go/internal/types"
)

func TestGetChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Data", "Sender", "Reciever", "Time"})

	message := []models.Message{
		{
			Id:       1,
			Data:     "data1",
			Sender:   1,
			Receiver: 2,
			Time:     time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, item := range message {
		rows = rows.AddRow(item.Id, item.Data, item.Sender,
			item.Receiver, item.Time)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, data, sender_id, 
		receiver_id, sent_time FROM message 
		WHERE ((sender_id = $1 AND receiver_id = $2) 
		OR (sender_id = $3 AND receiver_id = $4)) 
		ORDER BY sent_time`)).
		WithArgs(1, 2, 2, 1).
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	result, err := repo.GetChat(contexted, types.UserID(1), types.UserID(2))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(message, result) {
		t.Errorf("expected %v, got %v", message, result)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, data, sender_id, 
		receiver_id, sent_time FROM message 
		WHERE ((sender_id = $1 AND receiver_id = $2) 
		OR (sender_id = $3 AND receiver_id = $4)) 
		ORDER BY sent_time`)).
		WithArgs(1, 2, 2, 1).
		WillReturnError(fmt.Errorf("db_error"))

	result, err = repo.GetChat(contexted, types.UserID(1), types.UserID(2))
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if result != nil {
		t.Errorf("expected empty array, got %v", result)
		return
	}
}

func TestCreateMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{})

	tt := time.UnixMilli(01012001)

	mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO message 
		(data, sender_id, receiver_id, sent_time) 
		VALUES ($1,$2,$3,$4)`)).
		WithArgs("data1", 1, 2, tt).WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	message := models.MessageToReceive{
		Data:     "data1",
		Sender:   types.UserID(1),
		Receiver: types.UserID(2),
		Time:     01012001,
	}

	result, err := repo.CreateMessage(contexted, message)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(&message, result) {
		t.Errorf("expected %v, got %v", message, result)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO message 
		(data, sender_id, receiver_id, sent_time) 
		VALUES ($1,$2,$3,$4)`)).
		WithArgs("data1", 1, 2, tt).
		WillReturnError(fmt.Errorf("repo error"))

	result, err = repo.CreateMessage(contexted, message)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected no error, got %v", err)
	}

	if result != nil {
		t.Errorf("expected empty array, got %v", result)
		return
	}

}

func TestGetLastMessages(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Data", "Sender", "Reciever", "Time"})

	message := []models.Message{
		{
			Id:       1,
			Data:     "data1",
			Sender:   1,
			Receiver: 2,
			Time:     time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, item := range message {
		rows = rows.AddRow(item.Id, item.Data, item.Sender,
			item.Receiver, item.Time)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT DISTINCT 
		ON ( CASE WHEN sender_id < receiver_id 
			THEN sender_id || '_' || receiver_id ELSE receiver_id || '_' || 
			sender_id END) id, data, sender_id, receiver_id, sent_time 
			FROM message WHERE (sender_id = $1 OR receiver_id = $1) 
			AND ((sender_id = ANY($2)) OR (receiver_id = ANY($2))) 
			ORDER BY ( CASE WHEN sender_id < receiver_id 
				THEN sender_id || '_' || receiver_id 
				ELSE receiver_id || '_' || sender_id 
				END), sent_time DESC;`)).
		WithArgs(1, pq.Array([]int{2})).
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	result, err := repo.GetLastMessages(contexted, 1, []int{2})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(message, result) {
		t.Errorf("expected %v, got %v", message, result)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT DISTINCT 
		ON ( CASE WHEN sender_id < receiver_id 
			THEN sender_id || '_' || receiver_id ELSE receiver_id || '_' || 
			sender_id END) id, data, sender_id, receiver_id, sent_time 
			FROM message WHERE (sender_id = $1 OR receiver_id = $1) 
			AND ((sender_id = ANY($2)) OR (receiver_id = ANY($2))) 
			ORDER BY ( CASE WHEN sender_id < receiver_id 
				THEN sender_id || '_' || receiver_id 
				ELSE receiver_id || '_' || sender_id 
				END), sent_time DESC;`)).
		WithArgs(1, pq.Array([]int{2})).
		WillReturnError(fmt.Errorf("db_error"))

	result, err = repo.GetLastMessages(contexted, 1, []int{2})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if result != nil {
		t.Errorf("expected empty array, got %v", result)
		return
	}
}

func TestCreateClaim(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{})

	//tt := time.UnixMilli(01012001)

	mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO person_claim 
		(type, sender_id, receiver_id) 
		VALUES ($1,$2,$3)`)).
		WithArgs(1, 1, 2).
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	claim := models.Claim{
		Id:         1,
		TypeID:     1,
		SenderID:   1,
		ReceiverID: 2,
	}

	err = repo.CreateClaim(contexted, claim)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO person_claim 
		(type, sender_id, receiver_id) 
		VALUES ($1,$2,$3)`)).
		WithArgs(1, 1, 2).
		WillReturnError(fmt.Errorf("repo error"))

	err = repo.CreateClaim(contexted, claim)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected no error, got %v", err)
	}

}

func TestGetAllClaims(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Title"})

	message := []models.PureClaim{
		{
			Id:    1,
			Title: "title1",
		},
	}

	expMess := []models.PureClaim{
		{
			Id: 0,
		},
		{
			Id:    1,
			Title: "title1",
		},
	}

	for _, item := range message {
		rows = rows.AddRow(item.Id, item.Title)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, title FROM claim
		ORDER BY id`)).
		WithoutArgs().
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	result, err := repo.GetAllClaims(contexted)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(expMess, result) {
		t.Errorf("expected %v, got %v", expMess, result)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, title FROM claim
		ORDER BY id`)).
		WithoutArgs().
		WillReturnError(fmt.Errorf("db_error"))

	result, err = repo.GetAllClaims(contexted)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if result != nil {
		t.Errorf("expected empty array, got %v", result)
		return
	}
}

func TestGetImage(t *testing.T) {
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

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	images, err := repo.GetImages(contexted, 1)
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

	images, err = repo.GetImages(contexted, 1)
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

func TestGetInterest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name"})

	expect := []*models.Interest{
		{ID: 1, Name: "name1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT id, name FROM interest WHERE (id IN ($1))")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	filter := models.InterestGetFilter{
		ID: []types.InterestID{1},
		//Name: []string{"name1"},
	}

	images, err := repo.getInterests(contexted, &filter)
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
		regexp.QuoteMeta("SELECT id, name FROM interest WHERE (id IN ($1))")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	images, err = repo.getInterests(contexted, &filter)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if images != nil {
		t.Errorf("expected empty array, got %v", images)
		return
	}
}

func TestGetPersonInterest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	newRows := sqlmock.NewRows([]string{"personID", "interestID"})

	newExpect := []struct {
		interest types.InterestID
		person   types.UserID
	}{
		{interest: 1, person: 1},
	}

	for _, item := range newExpect {
		newRows = newRows.AddRow(item.interest, item.person)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT person_id, interest_id FROM person_interest WHERE person_id = $1")).
		WithArgs(1).
		WillReturnRows(newRows)

	rows := sqlmock.NewRows([]string{"ID", "Name"})

	expect := []*models.Interest{
		{ID: 1, Name: "name1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT id, name FROM interest WHERE (id IN ($1))")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	images, err := repo.GetPersonInterests(contexted, types.UserID(1))
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
		regexp.QuoteMeta("SELECT person_id, interest_id FROM person_interest WHERE person_id = $1")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	images, err = repo.GetPersonInterests(contexted, types.UserID(1))
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if images != nil {
		t.Errorf("expected empty array, got %v", images)
		return
	}
}

func TestGetLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Person1", "Person2"})

	message := []types.UserID{2}

	tmp := []models.Like{
		{
			Person1: types.UserID(1),
			Person2: types.UserID(2),
		},
	}

	for _, item := range tmp {
		rows = rows.AddRow(item.Person1, item.Person2)
	}

	// tmp := models.Like{
	// 	Person1: types.UserID(1),
	// 	Person2: types.UserID(2),
	// }

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT person_one_id, person_two_id 
		FROM "like" 
		WHERE (person_one_id = $1)`)).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	userID := types.UserID(1)

	filter := &models.LikeGetFilter{
		Person1: &userID,
	}

	result, err := repo.GetLike(contexted, filter)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(message, result) {
		t.Errorf("expected %v, got %v", message, result)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT person_one_id, person_two_id 
		FROM "like" 
		WHERE (person_one_id = $1)`)).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	result, err = repo.GetLike(contexted, filter)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if result != nil {
		t.Errorf("expected empty array, got %v", result)
		return
	}
}

func TestCreateLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{})

	//tt := time.UnixMilli(01012001)

	mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "like" 
		(person_one_id,person_two_id) 
		VALUES ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	err = repo.CreateLike(contexted, types.UserID(1), types.UserID(2))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "like" 
		(person_one_id,person_two_id) 
		VALUES ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnError(fmt.Errorf("repo error"))

	err = repo.CreateLike(contexted, types.UserID(1), types.UserID(2))
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected no error, got %v", err)
	}

}
