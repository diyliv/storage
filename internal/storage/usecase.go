package storage

import (
	"context"

	"github.com/diyliv/storage/internal/models"
)

type Usecase interface {
	Register(ctx context.Context, user models.User) error
	GetUserInfo(ctx context.Context, email string) (models.User, error)
	SavePublicKey(ctx context.Context, userId int, key, passPhrase string) error
	DeleteUserByEmail(ctx context.Context, email string) error
	CreateSession(ctx context.Context, userId, userName, userEmail, sessionToken string) error
	GetSessionInfo(ctx context.Context, sessionToken string) (map[string]string, error)
	CheckToken(ctx context.Context, sessionToken string) error
}
