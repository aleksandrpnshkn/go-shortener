package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	ServerAddr      string
	PublicBaseURL   string
	LogLevel        string
	FileStoragePath string
	DatabaseDSN     string
}

func New() *Config {
	fileStoragePath := ""
	tempDir := os.TempDir()

	if len(tempDir) != 0 {
		fileStoragePath = tempDir + "/go_shortener_storage.txt"
	}

	config := Config{
		ServerAddr:      "localhost:8080",
		PublicBaseURL:   "http://localhost:8080",
		LogLevel:        "info",
		FileStoragePath: fileStoragePath,
		DatabaseDSN:     "host=127.0.0.1 port=5432 user=admin password=qwerty dbname=shortener sslmode=disable",
	}

	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "net address host:port")
	flag.StringVar(&config.PublicBaseURL, "b", config.PublicBaseURL, "public base url for short links")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "file storage path")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "database DSN")

	flag.Parse()

	envServerAddr := os.Getenv("SERVER_ADDRESS")
	if envServerAddr != "" {
		config.ServerAddr = envServerAddr
	}

	envPublicBaseURL := os.Getenv("BASE_URL")
	if envPublicBaseURL != "" {
		config.PublicBaseURL = envPublicBaseURL
	}
	config.PublicBaseURL = strings.TrimRight(config.PublicBaseURL, "/")

	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel != "" {
		config.LogLevel = envLogLevel
	}

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		config.FileStoragePath = envFileStoragePath
	}

	envDatabaseDSN := os.Getenv("DATABASE_DSN")
	if envDatabaseDSN != "" {
		config.DatabaseDSN = envDatabaseDSN
	}

	return &config
}
