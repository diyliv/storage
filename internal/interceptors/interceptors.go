package interceptors

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diyliv/storage/internal/storage"
	"github.com/diyliv/storage/pkg/errs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type interceptor struct {
	logger    *zap.Logger
	storageUC storage.Usecase
}

func NewInterceptor(logger *zap.Logger, storageUC storage.Usecase) *interceptor {
	return &interceptor{
		logger:    logger,
		storageUC: storageUC,
	}
}

func (i *interceptor) Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	i.logger.Info(fmt.Sprintf("Method: %s, Time: %v, Metadata: %v, Err: %v\n",
		info.FullMethod,
		time.Since(start),
		md,
		err))
	return reply, err
}

func (i *interceptor) CheckToken(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	if info.FullMethod == "/storage.StorageService/CreateSession" || info.FullMethod == "/storage.StorageService/Register" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "unable to recieve medata: "+err.Error())
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "you need specify a token")
	}

	token := authHeader[0]
	if err := i.storageUC.CheckToken(ctx, token); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "this token is invalid")
		}
	}

	reply, err := handler(ctx, req)
	return reply, err
}
