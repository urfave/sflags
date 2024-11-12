package main

// This packages shows how to use sflags with kingpin library.

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/sflags"
	"github.com/urfave/sflags/gen/gkingpin"
)

type httpConfig struct {
	Host    string `desc:"HTTP host2"`
	Port    int    `flag:"port p"`
	SSL     bool
	Timeout time.Duration
	Addr    *net.TCPAddr
}

type config struct {
	HTTP       httpConfig
	Regexp     *regexp.Regexp
	Count      sflags.Counter
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

	app := kingpin.New("testApp", "")
	app.Terminate(nil)

	err := gkingpin.ParseTo(cfg, app, sflags.Prefix("kingpin."))
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	// print usage
	_, err = app.Parse([]string{"--help"})
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	_, err = app.Parse([]string{
		"--kingpin.http-host", "localhost",
		"-p", "9000",
		"--kingpin.http-ssl",
		"--kingpin.http-timeout", "30s",
		"--kingpin.http-addr", "google.com:8000",
		"--kingpin.regexp", "ddfd",
		"--kingpin.count", "--kingpin.count",
		"--kingpin.hidden-flag", "hidden_value",
	})
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Printf("\ncfg: %s\n", spew.Sdump(cfg))
}
