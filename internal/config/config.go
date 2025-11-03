package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerAddr      string
	PublicBaseURL   string
	LogLevel        string
	FileStoragePath string
	DatabaseDSN     string
	AuthSecretKey   string
	Audit           AuditConfig
	EnablePprof     bool
}

type AuditConfig struct {
	File string
	URL  string
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
		AuthSecretKey:   "changeme",
		EnablePprof:     true,
		Audit: AuditConfig{
			File: "",
			URL:  "",
		},
	}

	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "net address host:port")
	flag.StringVar(&config.PublicBaseURL, "b", config.PublicBaseURL, "public base url for short links")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "file storage path")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "database DSN")
	flag.StringVar(&config.AuthSecretKey, "s", config.DatabaseDSN, "auth secret key for signing JWT tokens")
	flag.StringVar(&config.Audit.File, "audit-file", config.Audit.File, "file path to store audit logs")
	flag.StringVar(&config.Audit.URL, "audit-url", config.Audit.URL, "external service URL to store audit logs")
	flag.BoolVar(&config.EnablePprof, "enable-pprof", config.EnablePprof, "enable pprof debug routes")

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

	authSecretKey, ok := os.LookupEnv("AUTH_SECRET_KEY")
	if ok {
		config.AuthSecretKey = authSecretKey
	}

	auditFile, ok := os.LookupEnv("AUDIT_FILE")
	if ok {
		config.Audit.File = auditFile
	}

	auditURL, ok := os.LookupEnv("AUDIT_URL")
	if ok {
		config.Audit.URL = auditURL
	}

	enablePprof, ok := os.LookupEnv("ENABLE_PPROF")
	if ok {
		enablePprofVar, err := strconv.ParseBool(enablePprof)
		if err != nil {
			enablePprofVar = false
		}
		config.EnablePprof = enablePprofVar
	}

	return &config
}
