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
		DatabaseDSN:     "postgres://admin:qwerty@localhost:5432/shortener?sslmode=disable",
	}

	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "net address host:port")
	flag.StringVar(&config.PublicBaseURL, "b", config.PublicBaseURL, "public base url for short links")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "file storage path")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "database DSN")

	flag.Parse()

	envServerAddr, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		config.ServerAddr = envServerAddr
	}

	envPublicBaseURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		config.PublicBaseURL = envPublicBaseURL
	}
	config.PublicBaseURL = strings.TrimRight(config.PublicBaseURL, "/")

	envLogLevel, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		config.LogLevel = envLogLevel
	}

	envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		config.FileStoragePath = envFileStoragePath
	}

	envDatabaseDSN, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		config.DatabaseDSN = envDatabaseDSN
	}

	return &config
}
