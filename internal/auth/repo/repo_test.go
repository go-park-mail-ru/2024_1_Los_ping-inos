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
	"github.com/go-redis/redismock/v9"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	models "main.go/internal/auth"
	. "main.go/internal/logs"
	"main.go/internal/types"
)

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

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	err = repo.CreatePersonInterests(contexted, 1, []types.InterestID{1})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO person_interest (person_id,interest_id) VALUES ($1,$2)")).
		WithArgs(1, 1).WillReturnError(fmt.Errorf("repo error"))

	err = repo.CreatePersonInterests(contexted, 1, []types.InterestID{1})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected no error, got %v", err)
	}

}

func TestActivateSub(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{})

	mock.ExpectQuery(regexp.QuoteMeta("UPDATE person SET premium = $1, premium_expires_at = $2 WHERE id = $3")).
		WithArgs(true, sqlmock.AnyArg(), 1).WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO person_payment (person_id,paymentTime) VALUES ($1,$2)")).
		WithArgs(1, sqlmock.AnyArg()).WillReturnRows(rows)

	repo := &PersonStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	err = repo.ActivateSub(contexted, 1, time.Now())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery(regexp.QuoteMeta("UPDATE person SET premium = $1, premium_expires_at = $2 WHERE id = $3")).
		WithArgs(true, sqlmock.AnyArg(), 1).WillReturnError(fmt.Errorf("repo error"))

	err = repo.ActivateSub(contexted, 1, time.Now())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected no error, got %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("UPDATE person SET premium = $1, premium_expires_at = $2 WHERE id = $3")).
		WithArgs(true, sqlmock.AnyArg(), 1).WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO person_payment (person_id,paymentTime) VALUES ($1,$2)")).
		WithArgs(1, sqlmock.AnyArg()).WillReturnError(fmt.Errorf("repo error"))

	err = repo.ActivateSub(contexted, 1, time.Now())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected no error, got %v", err)
	}

}

func TestGetSubHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Time"})

	var tmp time.Time

	expect := []*models.PaymentHistory{
		{Times: []models.HistoryRecord{
			{
				Time:  tmp.Unix(),
				Sum:   "2",
				Title: "Подписка",
			},
		}},
	}

	rows = rows.AddRow(tmp)

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT paymentTime FROM person_payment WHERE person_id = $1")).
		WithArgs(1).
		WillReturnRows(rows)

	repo := &PersonStorage{
		dbReader: db,
	}

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	images, err := repo.GetSubHistory(contexted, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(images, expect[0]) {
		t.Errorf("expected %v, got %v", expect[0], images)
		return
	}

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT paymentTime FROM person_payment WHERE person_id = $1")).
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	images, err = repo.GetSubHistory(contexted, 1)
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
		"Password", "CreatedAt", "Premium", "LikesLeft", "Gender", "Expires"})

	var tt time.Time

	expect := []*models.Person{
		{ID: 1, Name: "name1", Birthday: time.Now(), Description: "baa", Location: "maam", Email: "email1",
			Password: "123", CreatedAt: time.Now(), Premium: false, LikesLeft: 5, Gender: "male", PremiumExpires: tt.Unix()},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Name, item.Birthday, item.Description, item.Location, item.Email,
			item.Password, item.CreatedAt, item.Premium, item.LikesLeft, item.Gender, tt)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT id, name, birthday, description, location, 
			email, password, created_at, premium, likes_left, gender, 
			premium_expires_at FROM person WHERE ((id IN ($1)) AND 
			LOWER(name) LIKE $2)`)).
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

	// mock.ExpectQuery(
	// 	regexp.QuoteMeta(`SELECT id, name, birthday, description, location,
	// 		email, password, created_at, premium, likes_left, gender,
	// 		premium_expires_at FROM person WHERE ((id IN ($1)) AND
	// 		LOWER(name) LIKE $2)`)).
	// 	WithArgs(1, "%name1%").
	// 	WillReturnError(fmt.Errorf("db_error"))

	// persons, err = repo.Get(contexted, &filter)
	// if err := mock.ExpectationsWereMet(); err != nil {
	// 	t.Errorf("there were unfulfilled expectations: %s", err)
	// 	return
	// }

	// if err == nil {
	// 	t.Errorf("expected error, got nil")
	// 	return
	// }

	// if persons != nil {
	// 	t.Errorf("expected empty array, got %v", persons)
	// 	return
	// }
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{})

	var tt time.Time

	person := models.Person{
		ID:             1,
		Name:           "name1",
		Birthday:       time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		Description:    "baa",
		Location:       "maam",
		Email:          "email1",
		Password:       "123",
		CreatedAt:      time.Now(),
		Premium:        true,
		LikesLeft:      5,
		Gender:         "male",
		PremiumExpires: tt.Unix(),
	}

	setMap := make(map[string]interface{})

	setMap["premiumExpires"] = float64(person.PremiumExpires)

	mock.ExpectQuery(regexp.QuoteMeta(`
					UPDATE person 
					SET ID = $1, birthday = $2, description = $3, 
					email = $4, gender = $5, 
					name = $6, password = $7, 
					premium = $8, premium_expires_at = $9 
					WHERE id = $10`)).
		WithArgs(float64(1), "0001-01-01T00:00:00Z", person.Description,
			person.Email, person.Gender, person.Name, person.Password, person.Premium, time.Unix(int64(setMap["premiumExpires"].(float64)), 0), person.ID).
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
					name = $6, password = $7, 
					premium = $8, premium_expires_at = $9 
					WHERE id = $10`)).
		WithArgs(float64(1), "0001-01-01T00:00:00Z", person.Description,
			person.Email, person.Gender, person.Name, person.Password, person.Premium, time.Unix(int64(setMap["premiumExpires"].(float64)), 0), person.ID).
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

	selectRow := "INSERT INTO person(name, birthday, email, password, gender, premium_expires_at) VALUES ($1, $2, $3, $4, $5, $6)"

	logger := InitLog()
	logger.RequestID = int64(1)

	//actions := []driver.Value{"click", "event"}

	mock.ExpectExec(
		regexp.QuoteMeta(selectRow)).
		WithArgs("name1", "01012001", "mail.com", sqlmock.AnyArg(), "male", sqlmock.AnyArg()).
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
		WithArgs("name1", "01012001", "mail.com", sqlmock.AnyArg(), "male", sqlmock.AnyArg()).
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

// func TestNewInterestStorage(t *testing.T) {
// 	db, err := sql.Open("postgres", "user=your_user password=your_password dbname=your_db sslmode=disable")
// 	if err != nil {
// 		t.Fatalf("Failed to open database: %v", err)
// 	}
// 	defer db.Close()

// 	storage := NewInterestStorage(db)
// 	if storage == nil {
// 		t.Fatalf("Failed to create ImageStorage instance")
// 	}
// }

func TestNewInterestStorage(t *testing.T) {
	dbReader, err := sql.Open("postgres", "source_name=mock")
	require.NoError(t, err)
	defer dbReader.Close()

	storage := NewInterestStorage(dbReader)
	require.NotNil(t, storage)
	require.Equal(t, dbReader, storage.dbReader)
}

func TestNewSessionStorage(t *testing.T) {
	db := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	stor := NewSessionStorage(db)
	assert.NotNil(t, stor)
	assert.Equal(t, stor.db, db)
}

// func TestSessionStorage_GetBySID(t *testing.T) {
// 	ctx := context.Background()
// 	db := &mocks.RedisClient{}
// 	stor := &SessionStorage{db: db}

// 	// Test case 1: Successful retrieval of session
// 	db.On("Get", ctx, "session-id").Return("123", nil).Once()
// 	session, err := stor.GetBySID(ctx, "session-id")
// 	require.NoError(t, err)
// 	assert.Equal(t, types.UserID(123), session.UID)

// 	// Test case 2: Error when retrieving session
// 	db.On("Get", ctx, "session-id").Return("", errors.New("error")).Once()
// 	_, err = stor.GetBySID(ctx, "session-id")
// 	require.Error(t, err)

// 	// Test case 3: Invalid session ID
// 	db.On("Get", ctx, "invalid-session-id").Return("", nil).Once()
// 	_, err = stor.GetBySID(ctx, "invalid-session-id")
// 	require.Error(t, err)

// 	// Test case 4: Error when converting session ID to int
// 	db.On("Get", ctx, "session-id").Return("invalid", nil).Once()
// 	_, err = stor.GetBySID(ctx, "session-id")
// 	require.Error(t, err)

// 	// Test case 5: Successful deletion of session
// 	db.On("Del", ctx, "session-id").Return(0, nil).Once()
// 	err = stor.DeleteSession(ctx, "session-id")
// 	require.NoError(t, err)
// }

func TestGetBySID(t *testing.T) {
	db, mock := redismock.NewClientMock()

	SID := "mami mama mami"

	//session := &auth.Session{SID: SID}

	mock.ExpectGet(SID).RedisNil()

	repo := NewSessionStorage(db)

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	_, err := repo.GetBySID(contexted, SID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateSession(t *testing.T) {
	db, mock := redismock.NewClientMock()

	SID := "mami mama mami"

	//session := &auth.Session{SID: SID}

	mock.ExpectDel(SID).RedisNil()

	repo := NewSessionStorage(db)

	logger := InitLog()
	logger.RequestID = int64(1)
	contexted := context.WithValue(context.Background(), Logg, logger)

	err := repo.DeleteSession(contexted, SID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewAuthPersonStorage(t *testing.T) {
	dbReader, err := sql.Open("postgres", "source_name=mock")
	require.NoError(t, err)
	defer dbReader.Close()

	storage := NewAuthPersonStorage(dbReader)
	require.NotNil(t, storage)
	require.Equal(t, dbReader, storage.dbReader)
}
