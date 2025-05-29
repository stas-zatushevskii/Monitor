package main

import "flag"

var (
	address string
)

func ParseFlags() {
	flag.StringVar(&address, "a", "127.0.0.1:8080", "port")
	flag.Parse()
}
