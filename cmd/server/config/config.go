package config

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
)

type AuditData struct {
	FilePath string
	URL      string
}

type Config struct {
	Address         string
	LogLevel        string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DSN             string
	HashKey         string
	Audit           AuditData
	PrivateKey      string
}

func ParseEnvToInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func defaultDataFile() string {
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), "data.json")
	}
	if wd, err := os.Getwd(); err == nil {
		return filepath.Join(wd, "data.json")
	}
	return "data.json"
}

func (cfg *Config) ParseEnv() error {
	if v := os.Getenv("ADDRESS"); v != "" {
		cfg.Address = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := ParseEnvToInt(os.Getenv("STORE_INTERVAL")); v != 0 {
		cfg.StoreInterval = v
	}
	if v := os.Getenv("FILE_STORAGE_PATH"); v != "" {
		cfg.FileStoragePath = v
	}
	if v := os.Getenv("RESTORE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Restore = b
		}
	}
	if v := os.Getenv("DATABASE_DSN"); v != "" {
		cfg.DSN = v
	}
	if v := os.Getenv("KEY"); v != "" {
		cfg.HashKey = v
	}
	if v := os.Getenv("AUDIT_FILE"); v != "" {
		cfg.Audit.FilePath = v
	}
	if v := os.Getenv("AUDIT_URL"); v != "" {
		cfg.Audit.URL = v
	}
	if v := os.Getenv("CRYPTO_KEY"); v != "" {
		cfg.PrivateKey = v
	}
	return nil
}

func (cfg *Config) ParseFlags() {
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "port")
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "store interval in seconds")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultDataFile(), "file storage path")
	flag.BoolVar(&cfg.Restore, "r", false, "restore files")
	flag.StringVar(&cfg.DSN, "d", "", "database connection string") // postgres://user:pass@host:5432/db?sslmode=disable
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.StringVar(&cfg.Audit.FilePath, "audit-file", "", "path to logs file")
	flag.StringVar(&cfg.Audit.URL, "audit-url", "", "url where to send logs")
	flag.Parse()
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		FileStoragePath: defaultDataFile(),
	}
	cfg.ParseFlags()
	err := cfg.ParseEnv()
	return cfg, err
}
