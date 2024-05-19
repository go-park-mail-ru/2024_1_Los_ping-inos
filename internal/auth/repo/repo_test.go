package repo

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	models "main.go/internal/auth"
	. "main.go/internal/logs"
	"main.go/internal/types"
)

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

	rows := sqlmock.NewRows([]string{})

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

	repo := &InterestStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	filter := models.InterestGetFilter{
		ID: []types.InterestID{1},
		//Name: []string{"name1"},
	}

	images, err := repo.GetInterest(contexted, &filter)
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

	images, err = repo.GetInterest(contexted, &filter)
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

	repo := &InterestStorage{
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

func TestDeletePersonInterest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	selectRow := "DELETE FROM person_interest WHERE (person_id = $1 AND interest_id IN ($2))"

	logger := InitLog()
	logger.RequestID = int64(1)

	rows := sqlmock.NewRows([]string{})

	mock.ExpectQuery(
		regexp.QuoteMeta(selectRow)).
		WithArgs(1, 1).WillReturnRows(rows)

	repo := &InterestStorage{
		dbReader: db,
	}

	contexted := context.WithValue(context.Background(), Logg, logger)

	err = repo.DeletePersonInterests(contexted, types.UserID(1), []types.InterestID{1})
	if err != nil {
		t.Errorf("repo error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(selectRow)).
		WithArgs(1, 1).WillReturnError(fmt.Errorf("repo err"))

	err = repo.DeletePersonInterests(contexted, types.UserID(1), []types.InterestID{1})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("repo error: %s", err)
	}
}

func TestGetPerson(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "Birthday", "Description", "Location", "Email",
		"Password", "CreatedAt", "Premium", "LikesLeft", "Gender"})

	expect := []*models.Person{
		{ID: 1, Name: "name1", Birthday: time.Now(), Description: "baa", Location: "maam", Email: "email1",
			Password: "123", CreatedAt: time.Now(), Premium: false, LikesLeft: 5, Gender: "male"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name, item.Birthday, item.Description, item.Location, item.Email,
			item.Password, item.CreatedAt, item.Premium, item.LikesLeft, item.Gender)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, name, birthday, description, location, email, 
		password, created_at, premium, likes_left, gender 
		FROM person WHERE ((id IN ($1)) AND name LIKE $2)`)).
		WithArgs(1, "%name1%").
		WillReturnRows(rows)

	repo := &PersonStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	filter := models.PersonGetFilter{
		ID:   []types.UserID{1},
		Name: "name1",
	}

	persons, err := repo.Get(contexted, &filter)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(persons, expect) {
		t.Errorf("expected %v, got %v", expect, persons)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, name, birthday, description, location, email, 
		password, created_at, premium, likes_left, gender 
		FROM person WHERE ((id IN ($1)) AND name LIKE $2)`)).
		WithArgs(1, "%name1%").
		WillReturnError(fmt.Errorf("db_error"))

	persons, err = repo.Get(contexted, &filter)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if persons != nil {
		t.Errorf("expected empty array, got %v", persons)
		return
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{})

	person := models.Person{
		ID:          1,
		Name:        "name1",
		Birthday:    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "baa",
		Location:    "maam",
		Email:       "email1",
		Password:    "123",
		CreatedAt:   time.Now(),
		Premium:     false,
		LikesLeft:   5,
		Gender:      "male",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
	UPDATE person 
	SET ID = $1, birthday = $2, description = $3, 
	email = $4, gender = $5, 
	name = $6, password = $7 
	WHERE id = $8`)).
		WithArgs(float64(1), "0001-01-01T00:00:00Z", person.Description,
			person.Email, person.Gender, person.Name, person.Password, person.ID).
		WillReturnRows(rows)

	repo := &PersonStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	err = repo.Update(contexted, person)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
	UPDATE person 
	SET ID = $1, birthday = $2, description = $3, 
	email = $4, gender = $5, 
	name = $6, password = $7 
	WHERE id = $8`)).
		WithArgs(float64(1), "0001-01-01T00:00:00Z", person.Description,
			person.Email, person.Gender, person.Name, person.Password, person.ID).
		WillReturnError(fmt.Errorf("repo error"))

	err = repo.Update(contexted, person)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected no error, got %v", err)
	}

}

func TestDeletePerson(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	selectRow := "DELETE FROM person WHERE id = $1"

	logger := InitLog()
	logger.RequestID = int64(1)

	rows := sqlmock.NewRows([]string{})

	mock.ExpectQuery(
		regexp.QuoteMeta(selectRow)).
		WithArgs(1).WillReturnRows(rows)

	repo := &PersonStorage{
		dbReader: db,
	}

	contexted := context.WithValue(context.Background(), Logg, logger)

	err = repo.Delete(contexted, types.UserID(1))
	if err != nil {
		t.Errorf("repo error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(selectRow)).
		WithArgs(1).WillReturnError(fmt.Errorf("repo err"))

	err = repo.Delete(contexted, types.UserID(1))
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("repo error: %s", err)
	}
}

func TestAddAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	selectRow := "INSERT INTO person(name, birthday, email, password, gender) " +
		"VALUES ($1, $2, $3, $4, $5)"

	logger := InitLog()
	logger.RequestID = int64(1)

	//actions := []driver.Value{"click", "event"}

	mock.ExpectExec(
		regexp.QuoteMeta(selectRow)).
		WithArgs("name1", "01012001", "mail.com", sqlmock.AnyArg(), "male").
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := &PersonStorage{
		dbReader: db,
	}

	contexted := context.WithValue(context.Background(), Logg, logger)

	_, err = repo.AddAccount(contexted, "name1", "01012001", "male", "mail.com", "password")
	if err != nil {
		t.Errorf("repo error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectExec(
		regexp.QuoteMeta(selectRow)).
		WithArgs("name1", "01012001", "mail.com", sqlmock.AnyArg(), "male").
		WillReturnError(fmt.Errorf("repo error"))

	_, err = repo.AddAccount(contexted, "name1", "01012001", "male", "mail.com", "password")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("repo error: %s", err)
	}
}

func TestGetMatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID"})

	expect := []types.UserID{1}

	for _, item := range expect {
		rows = rows.AddRow(item)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT t1.person_two_id FROM "like"t1 
		INNER JOIN "like"t2 ON t1.person_one_id = t2.person_two_id 
		AND t1.person_two_id = t2.person_one_id 
		WHERE (t1.person_one_id = $1) 
		ORDER BY t1.person_two_id`)).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &PersonStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	persons, err := repo.GetMatch(contexted, types.UserID(1))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(persons, expect) {
		t.Errorf("expected %v, got %v", expect, persons)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT t1.person_two_id FROM "like"t1 
		INNER JOIN "like"t2 ON t1.person_one_id = t2.person_two_id 
		AND t1.person_two_id = t2.person_one_id 
		WHERE (t1.person_one_id = $1) 
		ORDER BY t1.person_two_id`)).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	persons, err = repo.GetMatch(contexted, types.UserID(1))
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if persons != nil {
		t.Errorf("expected empty array, got %v", persons)
		return
	}
}
