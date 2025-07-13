package main

import (
	"context"
	"log"
	"os"

	"github.com/aleksandrpnshkn/go-shortener/internal/app"
	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/logs"
	"github.com/aleksandrpnshkn/go-shortener/internal/store"
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

	fileStorage, err := store.NewFileStorage(config.FileStoragePath)
	if err != nil {
		logger.Fatal("failed to init app store", zap.Error(err))
	}

	err = app.Run(ctx, config, logger, fileStorage)
	if err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
