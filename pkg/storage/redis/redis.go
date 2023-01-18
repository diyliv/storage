package redis

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/diyliv/storage/config"
)

func ConnRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		MinIdleConns: cfg.Redis.MinIdleConn,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout),
	})

	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}

	return client
}
