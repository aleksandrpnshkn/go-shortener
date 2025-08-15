package main

import (
	"context"
	"log"
	"os"

	"github.com/aleksandrpnshkn/go-shortener/internal/app"
	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/logs"
	"github.com/aleksandrpnshkn/go-shortener/internal/store"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	err = migrateDatabase(config.DatabaseDSN)
	if err != nil {
		logger.Warn("failed to run migrations: %v", zap.Error(err))
	}

	fileStorage, err := store.NewFileStorage(config.FileStoragePath)
	if err != nil {
		logger.Fatal("failed to init app store", zap.Error(err))
	}

	err = app.Run(ctx, config, logger, fileStorage)
	if err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}

func migrateDatabase(databaseDSN string) error {
	m, err := migrate.New("file://migrations", databaseDSN)
	if err != nil {
		return err
	}
	return m.Up()
}
