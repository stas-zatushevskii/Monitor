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
	AddressGRPC     string
	LogLevel        string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DSN             string
	HashKey         string
	Audit           AuditData
	PrivateKey      string
	TrustedSubnet   string
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
	if address, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Address = address
	}
	if addressGFRPC, ok := os.LookupEnv("ADDRESSGRPC"); ok {
		cfg.AddressGRPC = addressGFRPC
	}
	if logLvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		cfg.LogLevel = logLvl
	}
	if storeInterval, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		val := ParseEnvToInt(storeInterval)
		cfg.StoreInterval = val
	}
	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = fileStoragePath
	}
	if restore, ok := os.LookupEnv("RESTORE"); ok {
		if b, err := strconv.ParseBool(restore); err == nil {
			cfg.Restore = b
		}
	}
	if dbDBS, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DSN = dbDBS
	}
	if key, ok := os.LookupEnv("KEY"); ok {
		cfg.HashKey = key
	}
	if auditFile, ok := os.LookupEnv("AUDIT_FILE"); ok {
		cfg.Audit.FilePath = auditFile
	}
	if auditURL, ok := os.LookupEnv("AUDIT_URL"); ok {
		cfg.Audit.URL = auditURL
	}
	if keyC, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		cfg.PrivateKey = keyC
	}
	if subNet, ok := os.LookupEnv("TRUSTED_SUBNET"); ok {
		cfg.TrustedSubnet = subNet
	}
	return nil
}

func (cfg *Config) ParseFlags() {
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "host:port")
	flag.StringVar(&cfg.AddressGRPC, "ag", ":3200", "host:port")
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "store interval in seconds")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultDataFile(), "file storage path")
	flag.BoolVar(&cfg.Restore, "r", false, "restore files")
	flag.StringVar(&cfg.DSN, "d", "", "database connection string") // postgres://user:pass@host:5432/db?sslmode=disable
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.StringVar(&cfg.Audit.FilePath, "audit-file", "", "path to logs file")
	flag.StringVar(&cfg.Audit.URL, "audit-url", "", "url where to send logs")
	flag.StringVar(&cfg.Audit.URL, "t", "", "trusted subnet")
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
