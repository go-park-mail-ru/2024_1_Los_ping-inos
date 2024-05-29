package repo

import (
	"context"
	"database/sql"
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
			MsgType: "message",
			Properties: models.MsgProperties{
				Id:       1,
				Data:     "data1",
				Sender:   1,
				Receiver: 2,
				Time:     time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, item := range message {
		rows = rows.AddRow(item.Properties.Id, item.Properties.Data, item.Properties.Sender,
			item.Properties.Receiver, item.Properties.Time)
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
		MsgType: "1",
		Properties: struct {
			Id       int64        "json:\"id\""
			Data     string       "json:\"data\""
			Sender   types.UserID "json:\"sender\""
			Receiver types.UserID "json:\"receiver\""
			Time     int64        "json:\"time\""
		}{
			Id:       1,
			Data:     "data1",
			Sender:   types.UserID(1),
			Receiver: types.UserID(2),
			Time:     01012001,
		},
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
			MsgType: "message",
			Properties: models.MsgProperties{
				Id:       1,
				Data:     "data1",
				Sender:   1,
				Receiver: 2,
				Time:     time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, item := range message {
		rows = rows.AddRow(item.Properties.Id, item.Properties.Data, item.Properties.Sender,
			item.Properties.Receiver, item.Properties.Time)
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

	newRows := sqlmock.NewRows([]string{"ID"})

	ids := []struct {
		Id int64
	}{
		{
			Id: 1,
		},
	}

	for _, item := range ids {
		newRows = newRows.AddRow(item.Id)
	}

	//tt := time.UnixMilli(01012001)

	mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "like" 
		(person_one_id,person_two_id) 
		VALUES ($1,$2)`)).
		WithArgs(1, 2).
		WillReturnRows(rows)

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT person_one_id FROM "like" 
		WHERE (person_one_id = $1 AND person_two_id = $2)`)).
		WithArgs(2, 1).
		WillReturnRows(newRows)

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

func TestGetClaimed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID"})

	message := []types.UserID{2}

	for _, item := range message {
		rows = rows.AddRow(item)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT receiver_id 
		FROM person_claim 
		WHERE sender_id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	userID := types.UserID(1)

	result, err := repo.GetClaimed(contexted, userID)
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
		regexp.QuoteMeta(`SELECT receiver_id 
		FROM person_claim 
		WHERE sender_id = $1`)).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	result, err = repo.GetClaimed(contexted, userID)
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

func TestGetFeed(t *testing.T) {
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

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT person_one_id, person_two_id
		FROM "like"
		WHERE (person_one_id = $1)`)).
		WithArgs(1).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"ID"})

	message = []types.UserID{2}

	for _, item := range message {
		rows = rows.AddRow(item)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT receiver_id 
			FROM person_claim 
			WHERE sender_id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"ID", "Name", "Birthday", "Description", "Location", "Email",
		"Password", "CreatedAt", "Premium", "LikesLeft", "Gender", "Expires"})

	expect := []*models.Person{
		{ID: 1, Name: "name1", Birthday: time.Now(), Description: "baa", Location: "maam", Email: "email1",
			Password: "123", CreatedAt: time.Now(), Premium: false, LikesLeft: 5, Gender: "male"},
	}

	var tt time.Time

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name, item.Birthday, item.Description, item.Location, item.Email,
			item.Password, item.CreatedAt, item.Premium, item.LikesLeft, item.Gender, tt)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, name, birthday, description, location, email, 
		password, created_at, premium, likes_left, gender, premium_expires_at 
		FROM person WHERE id NOT IN ($1,$2,$3`)).
		WithArgs(2, 1, 2).
		WillReturnRows(rows)

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	result, err := repo.GetFeed(contexted, types.UserID(1))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(expect, result) {
		t.Errorf("expected %v, got %v", message, result)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT person_one_id, person_two_id
		FROM "like"
		WHERE (person_one_id = $1)`)).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	result, err = repo.GetFeed(contexted, types.UserID(1))
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

	rows = sqlmock.NewRows([]string{"Person1", "Person2"})

	message = []types.UserID{2}

	tmp = []models.Like{
		{
			Person1: types.UserID(1),
			Person2: types.UserID(2),
		},
	}

	for _, item := range tmp {
		rows = rows.AddRow(item.Person1, item.Person2)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT person_one_id, person_two_id
		FROM "like"
		WHERE (person_one_id = $1)`)).
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT receiver_id
				FROM person_claim
				WHERE sender_id = $1`)).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	result, err = repo.GetFeed(contexted, types.UserID(1))
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

	// rows = sqlmock.NewRows([]string{"Person1", "Person2"})

	// message = []types.UserID{2}

	// tmp = []models.Like{
	// 	{
	// 		Person1: types.UserID(1),
	// 		Person2: types.UserID(2),
	// 	},
	// }

	// for _, item := range tmp {
	// 	rows = rows.AddRow(item.Person1, item.Person2)
	// }

	// mock.ExpectQuery(
	// 	regexp.QuoteMeta(`SELECT person_one_id, person_two_id
	// 	FROM "like"
	// 	WHERE (person_one_id = $1)`)).
	// 	WithArgs(1).
	// 	WillReturnRows(rows)

	// rows = sqlmock.NewRows([]string{"ID"})

	// message = []types.UserID{2}

	// for _, item := range message {
	// 	rows = rows.AddRow(item)
	// }

	// mock.ExpectQuery(
	// 	regexp.QuoteMeta(`SELECT receiver_id
	// 		FROM person_claim
	// 		WHERE sender_id = $1`)).
	// 	WithArgs(1).
	// 	WillReturnRows(rows)

	// mock.ExpectQuery(
	// 	regexp.QuoteMeta(`SELECT id, name, birthday, description,
	// 		location, email, password, created_at, premium, likes_left,
	// 		gender FROM person WHERE id NOT IN ($1,$2,$3)`)).
	// 	WithArgs(2, 1, 2).
	// 	WillReturnError(fmt.Errorf("db_error"))

	// result, err = repo.GetFeed(contexted, types.UserID(1))
	// if err := mock.ExpectationsWereMet(); err != nil {
	// 	t.Errorf("there were unfulfilled expectations: %s", err)
	// 	return
	// }

	// if err == nil {
	// 	t.Errorf("expected error, got nil")
	// 	return
	// }

	// if result != nil {
	// 	t.Errorf("expected empty array, got %v", result)
	// 	return
	// }
}

func TestNewWebSocStorage(t *testing.T) {
	storage := NewWebSocStorage()
	if storage == nil {
		t.Fatalf("Failed to create ImageStorage instance")
	}
}

// func TestWSStorage_AddConnection(t *testing.T) {
// 	storage := &WSStorage{}
// 	ctx := context.Background()
// 	//conn := &websocket.Conn{}
// 	uid := types.UserID(1)

// 	// Test case 4: Adding a connection with a nil connection
// 	err := storage.AddConnection(ctx, nil, uid)
// 	require.Error(t, err)
// }

// func TestGetConnection_NilUserID(t *testing.T) {
// 	logger := InitLog()
// 	logger.RequestID = int64(1)
// 	contexted := context.WithValue(context.Background(), Logg, logger)
// 	storage := &WSStorage{}
// 	_, ok := storage.GetConnection(contexted, 0)
// 	assert.False(t, ok, "GetConnection should return false when userID is nil")
// }

func TestNewPostgresStorage(t *testing.T) {
	db, err := sql.Open("postgres", "user=your_user password=your_password dbname=your_db sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	storage := NewPostgresStorage(db)
	if storage == nil {
		t.Fatalf("Failed to create ImageStorage instance")
	}
}

func TestDecreaseLikesCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Premiun", "likes"})

	rows = rows.AddRow(false, 5)

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT premium, likes_left 
		FROM person 
		WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE person 
		SET likes_left = $1 
		WHERE id = $2`)).
		WithArgs(4, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &PostgresStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	userID := types.UserID(1)

	result, err := repo.DecreaseLikesCount(contexted, userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(4, result) {
		t.Errorf("expected %v, got %v", 4, result)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT premium, likes_left 
		FROM person 
		WHERE id = $1`)).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	result, err = repo.DecreaseLikesCount(contexted, userID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

}
