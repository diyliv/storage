package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/diyliv/storage/config"
	_ "github.com/lib/pq"
)

func ConnPostgres(cfg *config.Config) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port =%s user=%s password=%s sslmode=disable", cfg.Postgres.Host,
		cfg.Postgres.Port, cfg.Postgres.Login, cfg.Postgres.Password)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err := conn.Ping(); err != nil {
		panic(err)
	}

	conn.SetConnMaxLifetime(time.Minute * time.Duration(cfg.Postgres.ConnMaxLifeTime))
	conn.SetMaxOpenConns(cfg.Postgres.MaxOpenConn)
	conn.SetMaxIdleConns(cfg.Postgres.MaxIdleConn)

	return conn
}
