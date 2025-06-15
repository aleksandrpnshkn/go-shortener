package main

import (
	"github.com/aleksandrpnshkn/go-shortener/internal/app"
	"github.com/aleksandrpnshkn/go-shortener/internal/config"
)

func main() {
	config := config.InitConfig()
	app.Run(config)
}
