package repository

import (
	"context"
	"database/sql"

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
	rows, err := p.psql.Exec("INSERT INTO users(user_name, user_email, user_hashed_password) VALUES ($1, $2, $3)",
		user.UserName,
		user.UserEmail,
		user.UserHashedPassword)
	if rows == nil {
		return errs.ErrAlreadyExists
	}

	if err != nil {
		p.logger.Error("Error while creating new user: " + err.Error())
		return err
	}

	return nil
}

func (p *postgresRepo) GetUserInfo(ctx context.Context, email string) (models.User, error) {
	var user models.User

	rows, err := p.psql.Query("SELECT user_id, user_name, user_hashed_password FROM users WHERE user_email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, errs.ErrNotFound
		}
		p.logger.Error("Error while getting info about user: " + err.Error())
		return models.User{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.UserName, &user.UserHashedPassword); err != nil {
			p.logger.Error("Error while scanning values: " + err.Error())
			return models.User{}, err
		}
	}

	return user, nil
}

func (p *postgresRepo) CreateSession(ctx context.Context, id int, sessionToken string) error {
	return nil
}

func (p *postgresRepo) SavePublicKey(ctx context.Context, key string) error {
	return nil
}
