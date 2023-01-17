package storage

import (
	"context"

	"github.com/diyliv/storage/internal/models"
)

type Usecase interface {
	Register(ctx context.Context, user models.User) error
	GetUserInfo(ctx context.Context, email string) (models.User, error)
	CreateSession(ctx context.Context, id int, sessionToken string) error
	SavePublicKey(ctx context.Context, key string) error
}
