package migration

import (
	"context"
	"errors"
	"fmt"
	"loverly/lib/log"
	"loverly/src/config"

	pg "loverly/lib/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrateLogIdentifier = "loverly"
)

type MigrationService interface {
	Up(context.Context) error
	Rollback(context.Context) error
	Version(context.Context) (int, bool, error)
}

type migrationService struct {
	driver  database.Driver
	migrate *migrate.Migrate
	log     log.Interface
}

func New(ctx context.Context, log log.Interface, cfg config.Postgres) (MigrationService, error) {
	pgCfg := pg.PostgresConfig{
		ConnectionUrl: cfg.ConnURI,
	}

	sqlxDB, err := pg.InitSQLX(ctx, log, pgCfg)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("error connecting to sqlxDB url:%s, err: %v", pgCfg.ConnectionUrl, err))
		return nil, err
	}

	databaseInstance, err := postgres.WithInstance(sqlxDB.DB, &postgres.Config{})
	if err != nil {
		log.Error(ctx, fmt.Sprintf("go-migrate postgres drv init failed: %v", err))
		return nil, err
	}

	migrate, err := migrate.NewWithDatabaseInstance("file://migration/sql",
		migrateLogIdentifier, databaseInstance)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("migrate init failed %v", err))
		return nil, err
	}

	return migrationService{
		driver:  databaseInstance,
		migrate: migrate,
		log:     log,
	}, nil
}

func (s migrationService) Up(ctx context.Context) error {
	currVersion, _, err := s.Version(ctx)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("Failed get current version err: %v", err))
		return err
	}

	s.log.Info(ctx, fmt.Sprintf("Running migration from version: %d", currVersion))
	if err := s.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			s.log.Info(ctx, "No Changes")
			return nil
		}
		s.log.Error(ctx, fmt.Sprintf("Failed run migrate err: %v", err))
		return err
	}

	currVersion, _, _ = s.Version(ctx)
	s.log.Error(ctx, fmt.Sprintf("Migration success, current version: %v", currVersion))
	return nil
}

func (s migrationService) Rollback(ctx context.Context) error {
	currVersion, _, err := s.Version(ctx)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("Failed get current version err: %v", err))
		return err
	}

	s.log.Info(ctx, fmt.Sprintf("Rollingback 1 step from version: %d", currVersion))

	if err := s.migrate.Steps(-1); err != nil {
		s.log.Error(ctx, fmt.Sprintf("Failed to rollback, err:%v", err))
		return err
	}

	currVersion, _, _ = s.Version(ctx)
	s.log.Info(ctx, fmt.Sprintf("Rollback success, current version:%d", currVersion))
	return nil
}

func (s migrationService) Version(ctx context.Context) (int, bool, error) {
	currVersion, dirty, err := s.driver.Version()
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("Failed to get version: %v", err))
		return 0, false, err
	}
	return currVersion, dirty, nil
}
