package main

// This packages shows how to use sflags with urfave/cli library.

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli/v2"
	"github.com/urfave/sflags"
	"github.com/urfave/sflags/gen/gcli"
)

type httpConfig struct {
	Host    string ` desc:"HTTP host"`
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

	flags, err := gcli.Parse(cfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	cliApp := cli.NewApp()
	cliApp.Action = func(c *cli.Context) error {
		return nil
	}
	cliApp.Flags = flags
	// print usage
	err = cliApp.Run([]string{"cliApp", "--help"})
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	err = cliApp.Run([]string{
		"cliApp",
		"--count=10",
		"--http-host", "localhost",
		"-p", "9000",
		"--http-ssl",
		"--http-timeout", "30s",
		"--http-addr", "google.com:8000",
		"--regexp", "ddfd",
		"--count", "--count",
		"--hidden-flag", "hidden_value",
	})
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Printf("\ncfg: %s\n", spew.Sdump(cfg))
}
