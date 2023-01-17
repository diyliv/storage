package server

import (
	"go.uber.org/zap"

	"github.com/diyliv/storage/config"
)

type server struct {
	logger *zap.Logger
	cfg    *config.Config
}

func NewServer(logger *zap.Logger, cfg *config.Config) *server {
	return &server{
		logger: logger,
		cfg:    cfg,
	}
}
func (s *server) StartgRPC() {
}
