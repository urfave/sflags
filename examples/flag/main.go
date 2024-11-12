package main

// This packages shows how to use sflags with flag library.

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/sflags"
	"github.com/urfave/sflags/gen/gflag"
)

type httpConfig struct {
	Host    string `desc:"HTTP host"`
	Port    int
	SSL     bool
	Timeout time.Duration
	Addr    *net.TCPAddr
	Methods map[string]int64 `desc:"HTTP Methods"`
}

type config struct {
	HTTP   httpConfig
	Regexp *regexp.Regexp
	Count  sflags.Counter
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
		Regexp: regexp.MustCompile("smth"),
		Count:  10,
	}
	// Use gflags.ParseToDef if you want default `flag.CommandLine`
	fs, err := gflag.Parse(cfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	// You should run fs.Parse(os.Args[1:]), but this is an example.
	err = fs.Parse([]string{
		"-count=20",
		"-http-host", "localhost",
		"-http-port", "9000",
		"-http-ssl",
		"-http-timeout", "30s",
		"-http-addr", "google.com:8000",
		"-regexp", "ddfd",
		"-count", "-count",
		"-http-methods", "post:15",
		"-http-methods", "get:25",
	})
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Println("Usage:")
	fs.SetOutput(os.Stdout)
	fs.PrintDefaults()
	fmt.Printf("cfg: %s\n", spew.Sdump(cfg))
}
