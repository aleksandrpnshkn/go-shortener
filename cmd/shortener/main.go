package main

import (
	"context"
	"log"
	"os"

	"github.com/aleksandrpnshkn/go-shortener/internal/app"
	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/logs"
	"github.com/aleksandrpnshkn/go-shortener/internal/store"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
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

	storage, err := store.NewStorage(ctx, config.DatabaseDSN, config.FileStoragePath, logger)
	if err != nil {
		logger.Fatal("failed to init app storage", zap.Error(err))
	}
	defer storage.Close()

	err = app.Run(ctx, config, logger, storage)
	if err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
