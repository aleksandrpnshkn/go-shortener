package main

import (
	"context"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"go.uber.org/zap"

	"github.com/aleksandrpnshkn/go-shortener/internal/app"
	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/logs"
	"github.com/aleksandrpnshkn/go-shortener/internal/store"
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

	urlsStorage, err := store.NewURLsStorage(ctx, config.DatabaseDSN, config.FileStoragePath, logger)
	if err != nil {
		logger.Fatal("failed to init app storage", zap.Error(err))
	}
	defer urlsStorage.Close()

	usersStorage, err := store.NewUsersStorage(ctx, config.DatabaseDSN, logger)
	if err != nil {
		logger.Fatal("failed to init users storage", zap.Error(err))
	}
	defer usersStorage.Close()

	err = app.Run(ctx, config, logger, urlsStorage, usersStorage)
	if err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
