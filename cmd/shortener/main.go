package main

import (
	"github.com/aleksandrpnshkn/go-shortener/internal/app"
	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/log"
	"github.com/aleksandrpnshkn/go-shortener/internal/store"
)

func main() {
	config := config.New()

	logger, err := log.NewLogger(config.LogLevel)
	if err != nil {
		panic("failed to create app logger")
	}
	defer logger.Sync()

	fileStore, err := store.NewFileStore(config.FileStoragePath)
	if err != nil {
		panic("failed to init app store")
	}

	app.Run(config, logger, fileStore)
}
