package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Address        string
	ReportInterval int
	PoolInterval   int
	HashKey        string
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
		cfg.PoolInterval = cfgE
	}
	if cfgE := os.Getenv("ADDRESS"); cfgE != "" {
		cfg.Address = cfgE
		return nil
	}
	if cfgE := os.Getenv("KEY"); cfgE != "" {
		cfg.HashKey = cfgE
	}
	return nil
}

func (cfg *Config) ParseFlags() {
	flag.IntVar(&cfg.ReportInterval, "r", 3, "report interval in seconds")
	flag.IntVar(&cfg.PoolInterval, "p", 2, "poll interval in seconds")
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "port")
	flag.StringVar(&cfg.HashKey, "k", "", "hash key")
	flag.Parse()
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	cfg.ParseFlags()
	err := cfg.ParseEnv()
	return cfg, err
}
