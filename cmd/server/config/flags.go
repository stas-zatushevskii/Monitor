package main

import (
	"flag"
	"os"
)

var address string

func ParseFlags() {
	flag.StringVar(&address, "a", "127.0.0.1:8080", "port")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		address = envRunAddr
	}

}
