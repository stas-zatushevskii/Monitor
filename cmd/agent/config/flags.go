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

	if reportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		val := ParseEnvToInt(reportInterval)
		cfg.ReportInterval = val
	}
	if pollInterval, ok := os.LookupEnv("POOL_INTERVAL"); ok {
		val := ParseEnvToInt(pollInterval)
		cfg.PollInterval = val
	}
	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.Address = addr
	}
	if addrGRPC, ok := os.LookupEnv("ADDRESSGRPC"); ok {
		cfg.AddressGRPC = addrGRPC
	}
	if key, ok := os.LookupEnv("KEY"); ok {
		cfg.HashKey = key
	}
	if rateLimit, ok := os.LookupEnv("RATE_LIMIT"); ok {
		val := ParseEnvToInt(rateLimit)
		cfg.RateLimit = val
	}
	if keyC, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		cfg.PublicKey = keyC
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
