package repository

import (
	"context"
	"testing"

	"github.com/diyliv/storage/config"
	"github.com/diyliv/storage/pkg/logger"
	"github.com/diyliv/storage/pkg/storage/redis"
)

var (
	cfg        = config.ReadConfig("../../../config")
	log        = logger.InitLogger()
	redisConn  = redis.ConnRedis(cfg)
	redisLogic = NewRedisRepo(log, redisConn, cfg)
	ctx        = context.Background()
)

func TestCreateSession(t *testing.T) {
	tc := []struct {
		userId       string
		userName     string
		userEmail    string
		sessionToken string
	}{
		{"1", "test_user_number_one", "test@email.com", "first_user_token"},
		{"2", "test_user_number_two", "second_test@email.com", "second_user_token"},
		{"3", "test_user_number_three", "third_test@email.com", "third_user_token"},
	}

	for _, val := range tc {
		if err := redisLogic.CreateSession(ctx, val.userId, val.userName, val.userEmail, val.sessionToken); err != nil {
			t.Errorf("Error while calling CreateSession(): %v\n", err)
		}
	}
}
