package usecase

import (
	"context"

	"github.com/diyliv/storage/internal/models"
	"github.com/diyliv/storage/internal/storage"
)

type storageUC struct {
	postgresRepo storage.PostgresRepository
	redisRepo    storage.RedisRepository
}

func NewStorageUC(postgresRepo storage.PostgresRepository, redisRepo storage.RedisRepository) *storageUC {
	return &storageUC{
		postgresRepo: postgresRepo,
		redisRepo:    redisRepo,
	}
}

func (s *storageUC) Register(ctx context.Context, user models.User) error {
	return s.postgresRepo.Register(ctx, user)
}

func (s *storageUC) GetUserInfo(ctx context.Context, email string) (models.User, error) {
	return s.postgresRepo.GetUserInfo(ctx, email)
}

func (s *storageUC) SavePublicKey(ctx context.Context, userId int, key, passPhrase string) error {
	return s.postgresRepo.SavePublicKey(ctx, userId, key, passPhrase)
}

func (s *storageUC) DeleteUserByEmail(ctx context.Context, email string) error {
	return s.postgresRepo.DeleteUserByEmail(ctx, email)
}

func (s *storageUC) CreateSession(ctx context.Context, userId, userName, userEmail, sessionToken string) error {
	return s.redisRepo.CreateSession(ctx, userId, userName, userEmail, sessionToken)
}

func (s *storageUC) CheckToken(ctx context.Context, sessionToken string) error {
	return s.redisRepo.CheckToken(ctx, sessionToken)
}

func (s *storageUC) GetSessionInfo(ctx context.Context, sessionToken string) (map[string]string, error) {
	return s.redisRepo.GetSessionInfo(ctx, sessionToken)
}
