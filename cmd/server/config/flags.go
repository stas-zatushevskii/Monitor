package config

import (
	"flag"
	"os"
)

var (
	Address      string
	FlagLogLevel string
)

func ParseFlags() {
	flag.StringVar(&Address, "a", "127.0.0.1:8080", "port")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		Address = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	}
}
