package main

// This packages shows how to use sflags with cobra library.
// cobra packages uses pflag.

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
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
	// set default values to config
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
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	err := gpflag.ParseTo(cfg, cmd.Flags())
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	// In you program use cmd.Execute(), but this is just an example.
	//if err := cmd.Execute(); err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}

	// this should show a message that old-flag is deprecated
	err = cmd.ParseFlags([]string{
		"--old-flag", "old",
	})
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	// parse normal command line arguments to see a changes in cfg structure
	err = cmd.ParseFlags([]string{
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
	cmd.Usage()
	//fmt.Printf("usage:\n%s\n", cmd.FlagUsages())
	fmt.Printf("\ncfg:\n %s", spew.Sdump(cfg))
}
