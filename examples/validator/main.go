package main

// This packages shows how to use sflags with flag library.

import (
	"log"

	"github.com/urfave/sflags"
	"github.com/urfave/sflags/gen/gflag"
	"github.com/urfave/sflags/validator/govalidator"
)

type config struct {
	Host string `valid:"host"`
	Port int    `valid:"port"`
}

func main() {
	cfg := &config{
		Host: "127.0.0.1",
		Port: 6000,
	}
	// Use gflags.ParseToDef if you want default `flag.CommandLine`
	fs, err := gflag.Parse(cfg, sflags.Validator(govalidator.New()))
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	// if we pass wrong domain to a host flag, we'll get a error.
	fs.Parse([]string{
		"-host", "wrong domain",
		"-port", "10",
	})
}
