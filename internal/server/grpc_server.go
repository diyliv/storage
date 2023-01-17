package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/diyliv/storage/config"
	"github.com/diyliv/storage/internal/storage"
	grpcservice "github.com/diyliv/storage/internal/storage/delivery/grpc"
	storagepb "github.com/diyliv/storage/proto/storage"
)

type server struct {
	logger    *zap.Logger
	cfg       *config.Config
	storageUC storage.Usecase
}

func NewServer(logger *zap.Logger, cfg *config.Config, storageUC storage.Usecase) *server {
	return &server{
		logger:    logger,
		cfg:       cfg,
		storageUC: storageUC,
	}
}
func (s *server) StartgRPC() {
	s.logger.Info(fmt.Sprintf("Starting gRPC server on port: %s\n", s.cfg.GrpcServer.Port))
	lis, err := net.Listen("tcp", s.cfg.GrpcServer.Port)
	if err != nil {
		s.logger.Error("Error while listening: " + err.Error())
	}
	defer lis.Close()

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: s.cfg.GrpcServer.MaxConnectionIdle * time.Minute,
			Timeout:           s.cfg.GrpcServer.Timeout * time.Second,
			MaxConnectionAge:  s.cfg.GrpcServer.MaxConnectionAge * time.Minute,
			Time:              s.cfg.GrpcServer.Timeout * time.Minute,
		}),
	}
	grpcservice := grpcservice.NewgRPCService(s.logger, s.storageUC)
	server := grpc.NewServer(opts...)
	storagepb.RegisterStorageServiceServer(server, grpcservice)

	go func() {
		if err := server.Serve(lis); err != nil {
			s.logger.Error("Error whiler serving: " + err.Error())
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done
	s.logger.Info("Exiting was successful")
	server.GracefulStop()
}
