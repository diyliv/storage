package main

import (
	"github.com/diyliv/storage/config"
	"github.com/diyliv/storage/internal/server"
	"github.com/diyliv/storage/internal/storage/repository"
	"github.com/diyliv/storage/internal/storage/usecase"
	"github.com/diyliv/storage/pkg/logger"
	"github.com/diyliv/storage/pkg/storage/postgres"
	"github.com/diyliv/storage/pkg/storage/redis"
)

func main() {
	cfg := config.ReadConfig("./config")
	logger := logger.InitLogger()

	redisConn := redis.ConnRedis(cfg)
	redisRepo := repository.NewRedisRepo(logger, redisConn, cfg)

	psqlConn, err := postgres.ConnPostgres(cfg)
	if err != nil {
		panic(err)
	}
	psqlRepo := repository.NewPostgresRepository(logger, psqlConn)
	psqlUC := usecase.NewStorageUC(psqlRepo, redisRepo)

	server := server.NewServer(logger, cfg, psqlUC)
	server.StartgRPC()
}
