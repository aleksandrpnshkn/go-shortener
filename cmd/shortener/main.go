package main

import (
	"github.com/aleksandrpnshkn/go-shortener/internal/app"
	"github.com/aleksandrpnshkn/go-shortener/internal/config"
	"github.com/aleksandrpnshkn/go-shortener/internal/log"
)

func main() {
	config := config.New()

	logger, err := log.NewLogger(config.LogLevel)
	if err != nil {
		panic("failed to create app logger")
	}
	defer logger.Sync()

	app.Run(config, logger)
}
