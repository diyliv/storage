package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"

	"github.com/diyliv/storage/config"
	"github.com/diyliv/storage/pkg/logger"
)

var (
	cfg       = config.ReadConfig("../../../config")
	log       = logger.InitLogger()
	ctx       = context.Background()
	redisConn = ConnRedis(&testing.T{})
	redisRep  = NewRedisRepo(log, redisConn, cfg)
)

func ConnRedis(t *testing.T) *redis.Client {
	mock, err := miniredis.Run()
	if err != nil {
		t.Errorf("Error while starting miniredis: %v\n", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         mock.Addr(),
		MinIdleConns: cfg.Redis.MinIdleConn,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout),
	})

	if err := client.Ping().Err(); err != nil {
		t.Errorf("Error while starting redis client: %v\n", err)
	}

	return client
}

func clear() {
	os.Remove("storage_service.json")
}

func TestCreateSession(t *testing.T) {
	defer clear()

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
		if err := redisRep.CreateSession(
			ctx,
			val.userId,
			val.userName,
			val.userEmail,
			val.sessionToken); err != nil {
			t.Errorf("Error while calling CreateSession() method: %v\n", err)
		}
	}
}

func TestCheckToken(t *testing.T) {
	defer clear()

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
		if err := redisRep.CheckToken(ctx, val.sessionToken); err != nil {
			t.Errorf("Error while checking token: %v\n", err)
		}
	}
}
