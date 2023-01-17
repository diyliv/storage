package grpc

import "go.uber.org/zap"

type grpcservice struct {
	logger *zap.Logger
}

func NewgRPCService(logger *zap.Logger) *grpcservice {
	return &grpcservice{
		logger: logger,
	}
}
