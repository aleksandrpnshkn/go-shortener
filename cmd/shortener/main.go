package main

import (
	"github.com/aleksandrpnshkn/go-shortener/internal/app"
	"github.com/aleksandrpnshkn/go-shortener/internal/config"
)

func main() {
	config := config.New()
	app.Run(config)
}
