package main

import (
	"context"
	"path/filepath"
	"runtime"

	"loverly/lib/i18n"
	"loverly/lib/jwt"
	"loverly/lib/log"
	"loverly/lib/postgres"
	"loverly/lib/redis"
	"loverly/src/business/domain"
	"loverly/src/business/usecase"
	"loverly/src/config"
	"loverly/src/handler"

	atomicSQLX "loverly/lib/atomic/sqlx"

	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()

	logger := log.Init(log.Config{Level: "Debug"})

	cfg, err := config.InitConfig(ctx, logger)
	if err != nil {
		panic(err)
	}

	if err := i18n.Init(ctx, cfg.Translation.FilePath, appTransFile, cfg.Translation.DefaultLanguage); err != nil {
		panic(err)
	}

	leader, err := postgres.InitSQLX(ctx, logger, postgres.PostgresConfig{
		ConnectionUrl:      cfg.Postgres.ConnURI,
		MaxPoolSize:        cfg.Postgres.MaxPoolSize,
		MaxIdleConnections: cfg.Postgres.MaxIdleConnections,
		ConnMaxIdleTime:    cfg.Postgres.MaxIdleTime,
		ConnMaxLifeTime:    cfg.Postgres.MaxLifeTime,
	})
	if err != nil {
		panic(err)
	}

	follower, err := postgres.InitSQLX(ctx, logger, postgres.PostgresConfig{
		ConnectionUrl:      cfg.PostgresReader.ConnURI,
		MaxPoolSize:        cfg.PostgresReader.MaxPoolSize,
		MaxIdleConnections: cfg.PostgresReader.MaxIdleConnections,
		ConnMaxIdleTime:    cfg.PostgresReader.MaxIdleTime,
		ConnMaxLifeTime:    cfg.PostgresReader.MaxLifeTime,
	})
	if err != nil {
		panic(err)
	}

	rds, err := redis.InitRedis(ctx, logger, cfg.Redis.Host, cfg.Redis.Password)
	if err != nil {
		panic(err)
	}

	dom := domain.Init(ctx, domain.InitParam{Log: logger, Cfg: cfg, LeaderDB: leader, FollowerDB: follower, Rds: rds})

	configJWT := jwt.Configuration{
		AccessTokenValidity:  cfg.AccessTokenValidity,
		RefreshTokenValidity: cfg.RefreshTokenValidity,
		IatLeeway:            cfg.IatLeeway,
		TokenIssuer:          cfg.TokenIssuer,
		KeyId:                cfg.JWTKey.KeyId,
		VerifyKey:            cfg.JWTKey.VerifyKey,
		SignKey:              cfg.JWTKey.SignKey,
	}

	jwt := jwt.Init(ctx, &configJWT, logger)

	tracer := otel.Tracer(cfg.ServiceName)

	atomicSessionProvider := atomicSQLX.NewSqlxAtomicSessionProvider(leader, tracer, logger)

	uc := usecase.Init(logger, *cfg, *jwt, *dom, atomicSessionProvider, tracer)

	handler.Init(ctx, logger, *cfg, uc, jwt)
}

var appTransFile = func() string {
	_, f, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(f))

	// Return the project root directory path.
	return filepath.Join(basepath, "translation")
}()
