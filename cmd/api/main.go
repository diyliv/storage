package main

import (
	"github.com/diyliv/storage/config"
	"github.com/diyliv/storage/internal/server"
	"github.com/diyliv/storage/internal/storage/repository"
	"github.com/diyliv/storage/internal/storage/usecase"
	"github.com/diyliv/storage/pkg/logger"
	"github.com/diyliv/storage/pkg/storage/postgres"
)

func main() {
	cfg := config.ReadConfig()
	logger := logger.InitLogger()

	psqlConn := postgres.ConnPostgres(cfg)
	psqlRepo := repository.NewPostgresRepository(logger, psqlConn)
	psqlUC := usecase.NewStorageUC(psqlRepo)

	server := server.NewServer(logger, cfg, psqlUC)
	server.StartgRPC()
}
