package interceptors

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/diyliv/storage/internal/storage"
	"github.com/diyliv/storage/pkg/errs"
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

	code, errStr := i.check(ctx)
	if code != codes.OK {
		return nil, status.Error(code, errStr)
	}

	reply, err := handler(ctx, req)
	return reply, err
}

func (i *interceptor) StreamLogger(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	start := time.Now()
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Error(codes.InvalidArgument, "unable to recieve medata")
	}

	err := handler(srv, stream)
	i.logger.Info(fmt.Sprintf("Method: %s, Time: %v, Metadata: %v, Error: %v\n",
		info.FullMethod,
		time.Since(start),
		md,
		err))
	return err
}

func (i *interceptor) StreamCheckToken(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	code, errStr := i.check(stream.Context())
	if code != codes.OK {
		return status.Error(code, errStr)
	}

	return handler(srv, stream)
}

func (i *interceptor) check(ctx context.Context) (codes.Code, string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return codes.InvalidArgument, "unable to recieve metadata"
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return codes.Unauthenticated, "you need to specify a token"
	}

	token := authHeader[0]
	if err := i.storageUC.CheckToken(ctx, token); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return codes.NotFound, "this token in invalid"
		}
		return codes.Internal, err.Error()
	}

	return codes.OK, ""
}
