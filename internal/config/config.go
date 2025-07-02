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
	}

	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "Net address host:port")
	flag.StringVar(&config.PublicBaseURL, "b", config.PublicBaseURL, "public base url for short links")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "file storage path")

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

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		config.FileStoragePath = envFileStoragePath
	}

	config.PublicBaseURL = strings.TrimRight(config.PublicBaseURL, "/")

	return &config
}
