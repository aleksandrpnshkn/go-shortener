package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	ServerAddr    string
	PublicBaseURL string
	LogLevel      string
}

func New() *Config {
	config := Config{
		ServerAddr:    "localhost:8080",
		PublicBaseURL: "http://localhost:8080",
		LogLevel:      "info",
	}

	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "Net address host:port")
	flag.StringVar(&config.PublicBaseURL, "b", config.PublicBaseURL, "public base url for short links")

	config.PublicBaseURL = strings.TrimRight(config.PublicBaseURL, "/")

	flag.Parse()

	envServerAddr := os.Getenv("SERVER_ADDRESS")
	if envServerAddr != "" {
		config.ServerAddr = envServerAddr
	}

	envPublicBaseURL := os.Getenv("BASE_URL")
	if envPublicBaseURL != "" {
		config.PublicBaseURL = envPublicBaseURL
	}

	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel != "" {
		config.LogLevel = envLogLevel
	}

	return &config
}
