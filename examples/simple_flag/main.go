package main

// This packages shows how to use sflags with flag library.
// Run this app with go run ./main.go -help

import (
	"flag"
	"log"
	"time"

	"github.com/urfave/sflags/gen/gflag"
)

type httpConfig struct {
	Host    string `desc:"HTTP host"`
	Port    int
	SSL     bool
	Timeout time.Duration
}

type config struct {
	HTTP httpConfig
}

func main() {
	cfg := &config{
		HTTP: httpConfig{
			Host:    "127.0.0.1",
			Port:    6000,
			SSL:     false,
			Timeout: 15 * time.Second,
		},
	}
	err := gflag.ParseToDef(cfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	flag.Parse()
}
