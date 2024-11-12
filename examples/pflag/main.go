package main

// This packages shows how to use sflags with pflag library.

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/sflags"
	"github.com/urfave/sflags/gen/gpflag"
)

type httpConfig struct {
	Host    string `desc:"HTTP host"`
	Port    int
	SSL     bool
	Timeout time.Duration
	Addr    *net.TCPAddr
}

type config struct {
	HTTP       httpConfig
	Regexp     *regexp.Regexp
	Count      sflags.Counter
	OldFlag    string `flag:",deprecated" desc:"use other flag instead"`
	HiddenFlag string `flag:",hidden"`
}

func main() {
	cfg := &config{
		HTTP: httpConfig{
			Host:    "127.0.0.1",
			Port:    6000,
			SSL:     false,
			Timeout: 15 * time.Second,
			Addr: &net.TCPAddr{
				IP:   net.ParseIP("127.0.0.1"),
				Port: 4000,
			},
		},
		Count:  12,
		Regexp: regexp.MustCompile("abc"),
	}

	fs, err := gpflag.Parse(cfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	err = fs.Parse([]string{
		"--old-flag", "old",
	})
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	err = fs.Parse([]string{
		"--count=10",
		"--http-host", "localhost",
		"--http-port", "9000",
		"--http-ssl",
		"--http-timeout", "30s",
		"--http-addr", "google.com:8000",
		"--regexp", "ddfd",
		"--count", "--count",
	})
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Printf("usage:\n%s\n", fs.FlagUsages())
	fmt.Printf("cfg: %s\n", spew.Sdump(cfg))
}
