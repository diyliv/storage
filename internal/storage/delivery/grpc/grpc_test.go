package grpc

import (
	"context"
	"database/sql"
	"net"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/diyliv/storage/config"
	"github.com/diyliv/storage/internal/storage/repository"
	"github.com/diyliv/storage/internal/storage/usecase"
	"github.com/diyliv/storage/pkg/logger"
	storagepb "github.com/diyliv/storage/proto/storage"
)

const bufSize = 1024 * 1024

var (
	cfg       = config.ReadConfig("../../../../config")
	log       = logger.InitLogger()
	db, _     = NewMock()
	sqlRepo   = repository.NewPostgresRepository(log, db)
	redisConn = ConnRedis(&testing.T{})
	redisRepo = repository.NewRedisRepo(log, redisConn, cfg)
	uc        = usecase.NewStorageUC(sqlRepo, redisRepo)
	grpcServ  = NewgRPCService(log, uc)
	lis       *bufconn.Listener
)

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	storagepb.RegisterStorageServiceServer(s, grpcServ)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

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

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Error("Error while creating new sqlmock db: " + err.Error())
	}

	return db, mock
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("Error while dialing with bufnet: %v\n", err)
	}
	defer conn.Close()

	client := storagepb.NewStorageServiceClient(conn)
	resp, err := client.Register(ctx, &storagepb.RegisterReq{
		UserName:     "some user",
		UserEmail:    "some@email.com",
		UserPassword: "hello world",
	})
	if err != nil {
		t.Errorf("Error while calling Register RPC: %v\n", err)
	}
	if resp.Status != "created" {
		t.Errorf("Unexpected result. Got %v want %v\n", resp.Status, "created")
	}
}

// func TestCreateSession(t *testing.T) {
// 	ctx := context.Background()
// 	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		t.Errorf("Error while dialing with bufnet: %v\n", err)
// 	}
// 	defer conn.Close()

// 	client := storagepb.NewStorageServiceClient(conn)
// 	_, err = client.CreateSession(ctx, &storagepb.CreateSessionReq{
// 		Email:    "first@email.com",
// 		Password: "hashed_pass",
// 	})

// 	if err != nil {
// 		t.Errorf("Error while calling CreateSession RPC: %v\n", err)
// 	}
// }
