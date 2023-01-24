package grpc

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/diyliv/storage/config"
	"github.com/diyliv/storage/internal/storage/repository"
	"github.com/diyliv/storage/internal/storage/usecase"
	"github.com/diyliv/storage/pkg/errs"
	"github.com/diyliv/storage/pkg/logger"
	"github.com/diyliv/storage/pkg/storage/postgres"
	storagepb "github.com/diyliv/storage/proto/storage"
)

const bufSize = 1024 * 1024

var (
	host     = "pg_test"
	port     = "5432"
	user     = "postgres"
	password = "postgres"
	db       = "postgres"
	cfg      = config.ReadConfig("../../../../config")
	log      = logger.InitLogger()
	conn, _  = postgres.ConnPostgres(&config.Config{Postgres: config.Postgres{
		Host:            host,
		Port:            port,
		Login:           user,
		Password:        password,
		DB:              db,
		ConnMaxLifeTime: 3,
		MaxOpenConn:     10,
		MaxIdleConn:     10,
	}})

	sqlRepo   = repository.NewPostgresRepository(log, conn)
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
		if errors.Is(err, errs.ErrAlreadyExists) {
			t.Logf("user already exists")
		}
		t.Errorf("Error while calling Register RPC: %v\n", err)
	}
	if resp.Status != "created" {
		t.Errorf("Unexpected result. Got %v want %v\n", resp.Status, "created")
	}
}

func TestCreateSession(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("Error while dialing with bufnet: %v\n", err)
	}
	defer conn.Close()
	client := storagepb.NewStorageServiceClient(conn)

	resp, err := client.Register(ctx, &storagepb.RegisterReq{
		UserName:     "some user",
		UserEmail:    "some_session@email.com",
		UserPassword: "hello world",
	})
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			t.Logf("user already exists")
		}
		t.Errorf("Error while calling Register RPC: %v\n", err)
	}
	if resp.Status != "created" {
		t.Errorf("Unexpected result. Got %v want %v\n", resp.Status, "created")
	}

	_, err = client.CreateSession(ctx, &storagepb.CreateSessionReq{
		Email:    "some_session@email.com",
		Password: "hello world",
	})

	if err != nil {
		t.Errorf("Error while calling CreateSession RPC: %v\n", err)
	}
}
