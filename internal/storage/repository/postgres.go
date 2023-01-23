package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/diyliv/storage/internal/models"
	"github.com/diyliv/storage/pkg/errs"
)

type postgresRepo struct {
	logger *zap.Logger
	psql   *sql.DB
}

func NewPostgresRepository(logger *zap.Logger, psql *sql.DB) *postgresRepo {
	return &postgresRepo{
		logger: logger,
		psql:   psql,
	}
}

func (p *postgresRepo) Register(ctx context.Context, user models.User) error {
	_, err := p.psql.Exec("INSERT INTO users (user_name, user_email, user_hashed_password) VALUES ($1, $2, $3)",
		user.UserName,
		user.UserEmail,
		user.UserHashedPassword)

	if err, ok := err.(*pq.Error); ok {
		if err.Code == pq.ErrorCode("23505") {
			return errs.ErrAlreadyExists
		}
		p.logger.Error("Error while creating new user: " + err.Error())
	}

	return nil
}

func (p *postgresRepo) GetUserInfo(ctx context.Context, email string) (models.User, error) {
	var user models.User

	rows, err := p.psql.Query("SELECT user_id, user_name, user_email, user_hashed_password FROM users WHERE user_email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errs.ErrNotFound
		}
		p.logger.Error("Error while getting info about user: " + err.Error())
		return models.User{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.UserName, &user.UserEmail, &user.UserHashedPassword); err != nil {
			p.logger.Error("Error while scanning values: " + err.Error())
			return models.User{}, err
		}
	}

	return user, nil
}

func (p *postgresRepo) SavePublicKey(ctx context.Context, userId int, key, passPharse string) error {
	_, err := p.psql.Exec("INSERT INTO users_keys (user_id, user_public_key, user_passphrase) VALUES ($1, $2, $3)",
		userId,
		key,
		passPharse)
	if err != nil {
		p.logger.Error("Error while saving public key: " + err.Error())
		return err
	}
	return nil
}

func (p *postgresRepo) DeleteUserByEmail(ctx context.Context, email string) error {
	_, err := p.psql.Query("DELETE FROM users WHERE user_email = $1", email)
	if err != nil {
		p.logger.Error("Error while deleting user: " + err.Error())
		return err
	}
	return nil
}

func (p *postgresRepo) Up() error {
	ctx := context.Background()

	query := "CREATE TABLE IF NOT EXISTS users (" +
		"user_id SERIAL PRIMARY KEY," +
		"user_name VARCHAR(16) NOT NULL," +
		"user_email VARCHAR(32) NOT NULL UNIQUE," +
		"user_hashed_password VARCHAR NOT NULL," +
		"user_updated_password TIMESTAMP DEFAULT NOW() NOT NULL," +
		"user_created_at TIMESTAMP DEFAULT NOW() NOT NULL);" +

		"CREATE TABLE IF NOT EXISTS users_keys(" +
		"user_id INT NOT NULL," +
		"user_public_key VARCHAR NOT NULL," +
		"user_passphrase VARCHAR(16) NOT NULL);"

	stmt, err := p.psql.PrepareContext(ctx, query)
	if err != nil {
		p.logger.Error("Error while creating tables: " + err.Error())
		return err
	}

	stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		p.logger.Error("Error while executing context: " + err.Error())
		return err
	}
	return nil
}

func (p *postgresRepo) Drop() error {
	ctx := context.Background()

	query := "DROP TABLE IF EXISTS users;" + "DROP TABLE IF EXISTS users_keys"

	stmt, err := p.psql.PrepareContext(ctx, query)
	if err != nil {
		p.logger.Error("Error while droppping tables: " + err.Error())
		return err
	}
	stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		p.logger.Error("Error while executing context: " + err.Error())
		return err
	}
	return nil
}

func (p *postgresRepo) Close() {
	p.psql.Close()
}
