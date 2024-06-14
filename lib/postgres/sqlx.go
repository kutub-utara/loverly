package postgres

import (
	"context"
	"fmt"

	"loverly/lib/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitSQLX(ctx context.Context, log log.Interface, cfg PostgresConfig) (*sqlx.DB, error) {
	xDB, err := otelsqlx.Open("postgres", cfg.ConnectionUrl, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to load the database err:%v", err))
		return nil, err
	}

	if err = xDB.Ping(); err != nil {
		log.Error(ctx, fmt.Sprintf("failed to ping the database err:%v", err))
		return nil, err
	}

	xDB.SetMaxOpenConns(cfg.MaxPoolSize)
	xDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	xDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	xDB.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	return xDB, nil
}
