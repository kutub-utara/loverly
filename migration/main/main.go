package main

import (
	"context"
	"fmt"
	"loverly/lib/log"
	"loverly/migration"
	"loverly/src/config"
	"os"
)

func main() {
	ctx := context.Background()

	logger := log.Init(log.Config{Level: "Debug"})

	cfg, err := config.InitConfig(ctx, logger)
	if err != nil {
		panic(err)
	}

	args := os.Args
	if len(args) < 2 {
		logger.Fatal(ctx, "Missing args. args: [up | rollback]")
	}

	migrationSvc, err := migration.New(ctx, logger, cfg.Postgres)
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Failed to initiate migration %v", err))
	}

	switch args[1] {
	case "up":
		migrationSvc.Up(ctx)
	case "rollback":
		migrationSvc.Rollback(ctx)
	default:
		logger.Fatal(ctx, "Invalid migration command")
	}
}
