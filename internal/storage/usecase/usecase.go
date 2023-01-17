package usecase

import (
	"context"

	"github.com/diyliv/storage/internal/models"
	"github.com/diyliv/storage/internal/storage"
)

type storageUC struct {
	postgresRepo storage.PostgresRepository
}

func NewStorageUC(postgresRepo storage.PostgresRepository) *storageUC {
	return &storageUC{postgresRepo: postgresRepo}
}

func (s *storageUC) Register(ctx context.Context, user models.User) error {
	return s.postgresRepo.Register(ctx, user)
}

func (s *storageUC) GetUserInfo(ctx context.Context, email string) (models.User, error) {
	return s.postgresRepo.GetUserInfo(ctx, email)
}

func (s *storageUC) CreateSession(ctx context.Context, id int, sessionToken string) error {
	return s.postgresRepo.CreateSession(ctx, id, sessionToken)
}

func (s *storageUC) SavePublicKey(ctx context.Context, key string) error {
	return s.postgresRepo.SavePublicKey(ctx, key)
}
