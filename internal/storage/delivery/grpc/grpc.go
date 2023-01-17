package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/diyliv/storage/internal/models"
	"github.com/diyliv/storage/internal/storage"
	"github.com/diyliv/storage/pkg/errs"
	"github.com/diyliv/storage/pkg/hash"
	rsaenc "github.com/diyliv/storage/pkg/rsa"
	storagepb "github.com/diyliv/storage/proto/storage"
)

type grpcservice struct {
	logger    *zap.Logger
	storageUC storage.Usecase
}

func NewgRPCService(logger *zap.Logger, storage storage.Usecase) *grpcservice {
	return &grpcservice{
		logger:    logger,
		storageUC: storage,
	}
}

func (gs *grpcservice) Register(ctx context.Context, req *storagepb.RegisterReq) (*storagepb.RegisterResp, error) {
	hashedPassword := hash.HashPass([]byte(req.GetUserPassword()))
	if err := gs.storageUC.Register(ctx, models.User{
		UserName:           req.GetUserName(),
		UserEmail:          req.GetUserEmail(),
		UserHashedPassword: hashedPassword,
	}); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "this user already exists")
		}
		gs.logger.Error("Error while calling Register() method: " + err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &storagepb.RegisterResp{Status: "created"}, status.Error(codes.OK, "'")
}

func (gs *grpcservice) CreateSession(ctx context.Context, req *storagepb.CreateSessionReq) (*storagepb.CreateSessionResp, error) {
	user, err := gs.storageUC.GetUserInfo(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "this user doesnt exist")
		}
		gs.logger.Error("Error while calling GetUserInfo() method: " + err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !hash.ComparePass(user.UserHashedPassword, []byte(req.Password)) {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// cache token in redis

	return &storagepb.CreateSessionResp{SessionToken: uuid.New().String()}, status.Error(codes.OK, "")
}

func (gs *grpcservice) Handshake(ctx context.Context, e *empty.Empty) (*storagepb.HandshakeResp, error) {
	keys, err := rsaenc.GenerateKeys()
	if err != nil {
		gs.logger.Error("Error while generating keys: " + err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	privateKey := fmt.Sprintf("%x", keys.D)

	return &storagepb.HandshakeResp{PrivateKey: privateKey}, status.Error(codes.OK, "")
}
