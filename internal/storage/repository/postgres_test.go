package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/diyliv/storage/internal/models"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Error("Error while creating new sqlmock db: " + err.Error())
	}

	return db, mock
}

func TestRegister(t *testing.T) {
	db, mock := NewMock()
	psql := NewPostgresRepository(log, db)
	defer db.Close()

	mock.ExpectExec("INSERT INTO users (user_name, user_email, user_hashed_password) VALUES ($1, $2, $3)").
		WithArgs("first_user", "first_user@email.com", "first_password").WillReturnResult(sqlmock.NewResult(1, 1))

	if err := psql.Register(ctx, models.User{
		UserName:           "first_user",
		UserEmail:          "first_user@email.com",
		UserHashedPassword: "first_password",
	}); err != nil {
		t.Errorf("Error while inserting user: %v\n", err)
	}

}

func TestGetUserInfo(t *testing.T) {
	db, mock := NewMock()
	psql := NewPostgresRepository(log, db)
	defer db.Close()

	mock.ExpectQuery("SELECT user_id, user_name, user_email, user_hashed_password FROM users WHERE user_email = $1").
		WithArgs("first@email.com").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "user_name", "user_email", "user_hashed_password"}).
			AddRow(1, "first_user", "first@email.com", "hashed_pass"))

	res, err := psql.GetUserInfo(ctx, "first@email.com")
	if err != nil {
		t.Errorf("Error while calling GetUserInfo(): %v\n", err)
	}

	if res.Id != 1 {
		t.Errorf("Unexpected value. Got %d want %d\n", res.Id, 1)
	}
	if res.UserName != "first_user" {
		t.Errorf("Unexpected value. Got %s want %s\n", res.UserName, "first_user")
	}
	if res.UserEmail != "first@email.com" {
		t.Errorf("Unexpected value. Got %s want %s\n", res.UserEmail, "first@email.com")
	}
	if res.UserHashedPassword != "hashed_pass" {
		t.Errorf("Unexpected value. Got %s want %s\n", res.UserHashedPassword, "hashed_pass")
	}
}

func TestSavePublicKey(t *testing.T) {
	db, mock := NewMock()
	psql := NewPostgresRepository(log, db)
	defer db.Close()

	mock.ExpectExec("INSERT INTO users_keys (user_id, user_public_key, user_passphrase) VALUES ($1, $2, $3)").
		WithArgs(1, "some_public_key", "some_passphrase").WillReturnResult(sqlmock.NewResult(1, 1))

	if err := psql.SavePublicKey(ctx, 1, "some_public_key", "some_passphrase"); err != nil {
		t.Errorf("Error while saving public key: %v\n", err)
	}
}

func TestDeleteUserByEmail(t *testing.T) {
	db, mock := NewMock()
	psql := NewPostgresRepository(log, db)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"user_id", "user_name", "user_email", "user_hashed_password", "user_updated_password", "user_created_at"}).
		AddRow(1, "first_user", "first@email.com", "hashed_pass", time.Now(), time.Now())

	mock.ExpectQuery("DELETE FROM users WHERE user_email = $1").
		WithArgs("first@email.com").WillReturnRows(rows)

	if err := psql.DeleteUserByEmail(ctx, "first@email.com"); err != nil {
		t.Errorf("Error while deleting user by email: %v\n", err)
	}
}
