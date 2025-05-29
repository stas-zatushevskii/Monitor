package main

import "flag"

var (
	reportIntervalFlag int
	pollIntervalFlag   int
	address            string
)

func ParseFlags() {
	flag.IntVar(&reportIntervalFlag, "r", 3, "report interval in seconds")
	flag.IntVar(&pollIntervalFlag, "p", 2, "poll interval in seconds")
	flag.StringVar(&address, "a", "127.0.0.1:8080", "port")
	flag.Parse()
}
