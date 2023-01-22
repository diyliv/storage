package repository

import (
	"database/sql"
	"testing"

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
	defer db.Close()

	mock.ExpectExec("INSERT INTO users (user_name, user_email, user_hashed_password) VALUES ($1, $2, $3)").
		WithArgs("first_user", "first_user@email.com", "first_password").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	psql := NewPostgresRepository(log, db)
	if err := psql.Register(ctx, models.User{
		UserName:           "first_user",
		UserEmail:          "first_user@email.com",
		UserHashedPassword: "first_password",
	}); err != nil {
		t.Errorf("Error while inserting user: %v\n", err)
	}

}
