package repository

import (
	"database/sql"
	"testing"

	"github.com/diyliv/storage/config"
	"github.com/diyliv/storage/internal/models"
	"github.com/diyliv/storage/pkg/storage/postgres"
	"github.com/lib/pq"
)

var (
	host     = "pg_test"
	port     = "5432"
	user     = "postgres"
	password = "postgres"
	db       = "postgres"
	conn, _  = postgres.ConnPostgres(&config.Config{Postgres: config.Postgres{
		Host:            host,
		Port:            port,
		Login:           user,
		Password:        password,
		DB:              db,
		ConnMaxLifeTime: 3,
		MaxOpenConn:     10,
		MaxIdleConn:     10,
	}})
	repo = NewPostgresRepository(log, conn)
)

func TestRegister(t *testing.T) {
	u := models.User{
		UserName:           "test",
		UserEmail:          "test@email.com",
		UserHashedPassword: "hashed_pass",
	}

	err := repo.Register(ctx, u)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == pq.ErrorCode("23505") {
			t.Logf("User already exists")
		}
		t.Errorf("Error while calling Register() db method: %v\n", err)
	}
}

func TestGetUserInfo(t *testing.T) {
	u := models.User{
		UserName:           "test",
		UserEmail:          "test_get_user_info@email.com",
		UserHashedPassword: "hashed_pass",
	}

	err := repo.Register(ctx, u)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == pq.ErrorCode("23505") {
			t.Logf("User already exists")
		}
		t.Errorf("Error while calling Register() db method: %v\n", err)
	}

	res, err := repo.GetUserInfo(ctx, u.UserEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			t.Errorf("Could not find this user: %v\n", err)
		}
		t.Errorf("Error while getting info about user: %v\n", err)
	}

	if res.UserName != u.UserName {
		t.Errorf("Unexpected values. Got %v want %v\n", res.UserName, u.UserName)
	}
	if res.UserHashedPassword != u.UserHashedPassword {
		t.Errorf("Unexpected values. Got %v want %v\n", res.UserHashedPassword, u.UserHashedPassword)
	}
}

func TestSavePublicKey(t *testing.T) {
	u := models.User{
		UserName:           "test",
		UserEmail:          "test_save_public_key@email.com",
		UserHashedPassword: "hashed_pass",
	}

	err := repo.Register(ctx, u)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == pq.ErrorCode("23505") {
			t.Logf("User already exists")
		}
		t.Errorf("Error while calling Register() db method: %v\n", err)
	}
	err = repo.SavePublicKey(ctx, 1, "some_key", "some_passphrase")
	if err != nil {
		t.Errorf("Error while saving public key: %v\n", err)
	}
}

func TestDeleteUserByEmail(t *testing.T) {
	u := models.User{
		UserName:           "test",
		UserEmail:          "test_delete_user_by@email.com",
		UserHashedPassword: "hashed_pass",
	}

	err := repo.Register(ctx, u)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == pq.ErrorCode("23505") {
			t.Logf("User already exists")
		}
		t.Errorf("Error while calling Register() db method: %v\n", err)
	}
	err = repo.DeleteUserByEmail(ctx, u.UserEmail)
	if err != nil {
		t.Errorf("Error while calling DeleteUserByEmail() db method: %v\n", err)
	}
}
