package repository

import (
	"context"

	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/diyliv/storage/config"
	"github.com/diyliv/storage/pkg/errs"
)

type redisRepo struct {
	logger *zap.Logger
	redis  *redis.Client
	cfg    *config.Config
}

func NewRedisRepo(logger *zap.Logger, redis *redis.Client, cfg *config.Config) *redisRepo {
	return &redisRepo{
		logger: logger,
		redis:  redis,
		cfg:    cfg,
	}
}

func (r *redisRepo) CreateSession(ctx context.Context, userId, userName, userEmail, sessionToken string) error {
	if err := r.redis.HMSet(sessionToken, map[string]interface{}{
		"userId":       userId,
		"userName":     userName,
		"userEmail":    userEmail,
		"sessionToken": sessionToken,
	}).Err(); err != nil {
		r.logger.Error("Error while caching token: " + err.Error())
		return err
	}

	return nil
}

func (r *redisRepo) CheckToken(ctx context.Context, sessionToken string) error {
	cmd := r.redis.HGet(sessionToken, "sessionToken")
	if cmd.Val() == "" {
		return errs.ErrNotFound
	}

	return nil
}
