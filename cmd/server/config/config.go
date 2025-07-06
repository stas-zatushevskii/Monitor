package config

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
)

var (
	Address         string
	FlagLogLevel    string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DSN             string
)

func ParseFlags() {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	DefaultFile := filepath.Join(filepath.Dir(exePath), "data.json")
	flag.StringVar(&Address, "a", "127.0.0.1:8080", "port")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.IntVar(&StoreInterval, "i", 300, "store interval in seconds")
	flag.StringVar(&FileStoragePath, "f", DefaultFile, "log level")
	flag.BoolVar(&Restore, "r", false, "restore files")
	flag.StringVar(&DSN, "d", "postgres://postgres:123@localhost:5432/postgres?sslmode=disable", "database connection string") // postgres://postgres:123@localhost:5432/postgres?sslmode=disable

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		Address = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		StoreInterval, err = strconv.Atoi(envStoreInterval)
		if err != nil {
			panic(err)
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		Restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			panic(err)
		}
	}
	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		DSN = envDSN
	}
}
