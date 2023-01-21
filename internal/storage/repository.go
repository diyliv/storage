package storage

import (
	"context"

	"github.com/diyliv/storage/internal/models"
)

type PostgresRepository interface {
	Register(ctx context.Context, user models.User) error
	GetUserInfo(ctx context.Context, email string) (models.User, error)
	SavePublicKey(ctx context.Context, userId int, key, passPhrase string) error
}

type RedisRepository interface {
	CreateSession(ctx context.Context, userId, userName, userEmail, sessionToken string) error
	CheckToken(ctx context.Context, sessionToken string) error
	GetSessionInfo(ctx context.Context, sessionToken string) (map[string]string, error)
}
