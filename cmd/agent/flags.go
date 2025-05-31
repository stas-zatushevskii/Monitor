package main

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Address        string
	ReportInterval int
	PoolInterval   int
}

func (cfg *Config) ParseEnv() error {

	if cfgE := os.Getenv("REPORT_INTERVAL"); cfgE != "" {
		cfgReportInterval, err := strconv.Atoi(cfgE)
		if err != nil {
			return err
		}
		cfg.ReportInterval = cfgReportInterval
	}
	if cfgE := os.Getenv("POOl_INTERVAL"); cfgE != "" {
		cfgReportInterval, err := strconv.Atoi(cfgE)
		if err != nil {
			return err
		}
		cfg.PoolInterval = cfgReportInterval
	}
	if cfgE := os.Getenv("ADDRESS"); cfgE != "" {
		cfg.Address = cfgE
		return nil
	}
	return nil
}

func (cfg *Config) ParseFlags() {
	flag.IntVar(&cfg.ReportInterval, "r", 3, "report interval in seconds")
	flag.IntVar(&cfg.PoolInterval, "p", 2, "poll interval in seconds")
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "port")
	flag.Parse()
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	cfg.ParseFlags()
	err := cfg.ParseEnv()
	return cfg, err
}
