package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"os"

	"github.com/aleksandrpnshkn/go-shortener/internal/app"
	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/logs"
	"github.com/aleksandrpnshkn/go-shortener/internal/store"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

func main() {
	config := config.New()
	ctx := context.Background()

	logger, err := logs.NewLogger(config.LogLevel)
	if err != nil {
		log.Printf("failed to create app logger: %v", err)
		os.Exit(1)
	}
	defer logger.Sync()

	var storage store.Storage

	sqlStorage, err := store.NewSQLStorage(ctx, config.DatabaseDSN)
	if err == nil {
		defer sqlStorage.Close()

		err = runMigrations(config.DatabaseDSN)
		if err != nil {
			logger.Fatal("failed to run migrations: %v", zap.Error(err))
		}

		storage = sqlStorage
	} else {
		logger.Warn("failed to init SQL storage", zap.Error(err))

		fileStorage, err := store.NewFileStorage(config.FileStoragePath)
		if err != nil {
			logger.Fatal("failed to init app storage", zap.Error(err))
		}

		storage = fileStorage
	}

	err = app.Run(ctx, config, logger, storage)
	if err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}

//go:embed migrations/*.sql
var migrationsFiles embed.FS

func runMigrations(databaseDSN string) error {
	sourceDriver, err := iofs.New(migrationsFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseDSN)
	if err != nil {
		return err
	}

	err = m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
