package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Address        string
	AddressGRPC    string
	ReportInterval int
	PollInterval   int
	HashKey        string
	RateLimit      int
	PublicKey      string
}

func ParseEnvToInt(cfg string) int {
	cfgDataInt, err := strconv.Atoi(cfg)
	if err != nil {
		return 0
	}
	return cfgDataInt
}

func (cfg *Config) ParseEnv() error {

	if cfgE := ParseEnvToInt(os.Getenv("REPORT_INTERVAL")); cfgE != 0 {
		cfg.ReportInterval = cfgE
	}
	if cfgE := ParseEnvToInt(os.Getenv("POOL_INTERVAL")); cfgE != 0 {
		cfg.PollInterval = cfgE
	}
	if cfgE := os.Getenv("ADDRESS"); cfgE != "" {
		cfg.Address = cfgE
	}
	if cfgE := os.Getenv("ADDRESSGRPC"); cfgE != "" {
		cfg.AddressGRPC = cfgE
	}
	if cfgE := os.Getenv("KEY"); cfgE != "" {
		cfg.HashKey = cfgE
	}
	if cfgE := ParseEnvToInt(os.Getenv("RATE_LIMIT")); cfgE != 0 {
		cfg.RateLimit = cfgE
	}
	if v := os.Getenv("CRYPTO_KEY"); v != "" {
		cfg.PublicKey = v
	}
	return nil
}

func (cfg *Config) ParseFlags() {
	flag.IntVar(&cfg.ReportInterval, "r", 3, "report interval in seconds")
	flag.IntVar(&cfg.PollInterval, "p", 2, "pool interval in seconds")
	flag.IntVar(&cfg.RateLimit, "l", 1, "rate limit")
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "host:port")
	flag.StringVar(&cfg.AddressGRPC, "ag", ":3200", "host:port")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.Parse()
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	cfg.ParseFlags()
	err := cfg.ParseEnv()
	return cfg, err
}
