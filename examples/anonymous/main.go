package main

// This packages shows how to use sflags with flag library.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/urfave/sflags/gen/gflag"
)

type httpConfig struct {
	Host string
}

type config struct {
	httpConfig
}

func main() {
	cfg := &config{
		httpConfig: httpConfig{
			Host: "127.0.0.1",
		},
	}
	// Use gflags.ParseToDef if you want default `flag.CommandLine`
	fs, err := gflag.Parse(cfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	fs.Init("", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	// You should run fs.Parse(os.Args[1:]), but this is an example.
	err = fs.Parse([]string{
		"-http-config-host", "localhost",
	})
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("Usage:")
	fs.SetOutput(os.Stdout)
	fs.PrintDefaults()
	fmt.Printf("cfg: %s\n", spew.Sdump(cfg))
}
