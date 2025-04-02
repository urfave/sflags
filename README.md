# Flags based on structures [![GoDoc](https://godoc.org/github.com/urfave/sflags?status.svg)](http://godoc.org/github.com/urfave/sflags)  [![Build Status](https://img.shields.io/github/check-runs/urfave/sflags/main?label=build%20status)](https://github.com/urfave/sflags/actions?query=branch%3Amain)  [![codecov](https://codecov.io/gh/urfave/sflags/branch/master/graph/badge.svg)](https://codecov.io/gh/urfave/sflags)  [![Go Report Card](https://goreportcard.com/badge/github.com/urfave/sflags)](https://goreportcard.com/report/github.com/urfave/sflags)

The sflags package uses structs, reflection and struct field tags
to allow you specify command line options. It supports [different types](#supported-types-in-structures) and [features](#features).

An example:

```golang
type HTTPConfig struct {
	Host    string        `desc:"HTTP host"`
	Port    int           `flag:"port p" desc:"some port"`
	SSL     bool          `env:"HTTP_SSL_VALUE"`
	Timeout time.Duration `flag:",deprecated,hidden"`
}

type Config struct {
	HTTP  HTTPConfig
	Stats StatsConfig
}
```

And you can use your favorite flag or cli library!

## Supported libraries and features:

|     |     | Hidden | Deprecated | Short | Env | Required |
| --- | --- |:------:|:----------:|:-----:|:---:|:--------:|
| <ul><li>[x] [flag]</li><ul> | [example](./examples/flag/main.go) | `-` | `-` | `-` | `-` | `-` |
| <ul><li>[x] [kingpin]</li></ul> | [example](./examples/kingpin/main.go) | <ul><li>[x] </li></ul> | <ul><li>[ ] </li></ul> | <ul><li>[x] </li></ul> | <ul><li>[x] </li></ul> | <ul><li>[x] </li></ul> |
| <ul><li>[x] [spf13/pflag]</li></ul> | [example](./examples/pflag/main.go) | <ul><li>[x] </li></ul> | <ul><li>[x] </li></ul> | <ul><li>[x] </li></ul> | `-` | `-` |
| <ul><li>[x] [spf13/cobra]</li></ul> | [example](./examples/cobra/main.go) | <ul><li>[x] </li></ul> | <ul><li>[x] </li></ul> | <ul><li>[x] </li></ul> | `-` | `-` |
| <ul><li>[x] [urfave/cli]</li></ul> | [example](./examples/urfave_cli/main.go) | <ul><li>[x] </li></ul> | `-` | <ul><li>[x] </li></ul> | <ul><li>[x] </li></ul> | <ul><li>[x] </li></ul> |

- [x] - feature is supported and implemented

`-` - feature can't be implemented for this cli library


[flag]: https://golang.org/pkg/flag/
[spf13/pflag]: https://github.com/spf13/pflag
[spf13/cobra]: https://github.com/spf13/cobra
[spf13/viper]: https://github.com/spf13/viper
[urfave/cli]: https://github.com/urfave/cli
[kingpin]: https://github.com/alecthomas/kingpin

## Features:

 - [x] Set environment name
 - [x] Set usage
 - [x] Long and short forms
 - [x] Skip field
 - [x] Required
 - [ ] Placeholders (by `name`)
 - [x] Deprecated and hidden options
 - [x] Multiple ENV names
 - [x] Interface for user types.
 - [x] [Validation](https://godoc.org/github.com/urfave/sflags/validator/govalidator#New) (using [govalidator](https://github.com/asaskevich/govalidator) package)
 - [x] Anonymous nested structure support (anonymous structures flatten by default)

## Supported types in structures:

 - [x] `int`, `int8`, `int16`, `int32`, `int64`
 - [x] `uint`, `uint8`, `uint16`, `uint32`, `uint64`
 - [x] `float32`, `float64`
 - [x] slices for all previous numeric types (e.g. `[]int`, `[]float64`)
 - [x] `bool`
 - [x] `[]bool`
 - [x] `string`
 - [x] `[]string`
 - [x] nested structures
 - [x] net.TCPAddr
 - [x] net.IP
 - [x] time.Duration
 - [x] regexp.Regexp
 - [x] map for all previous types (e.g. `map[int64]bool`, `map[string]float64`)

## Custom types:
 - [x] HexBytes

 - [x] count
 - [ ] ipmask
 - [ ] enum values
 - [ ] enum list values
 - [ ] file
 - [ ] file list
 - [ ] url
 - [ ] url list
 - [ ] units (bytes 1kb = 1024b, speed, etc)

## Example:

The code below shows how to use `sflags` with the [flag] library. Examples for other
flag libraries are available from [./examples](./examples) dir.

```golang
package main

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
```

That code generates next output:
```sh
$ go run ./main.go --help
Usage of _obj/exe/main:
  -http-host value
    	HTTP host (default 127.0.0.1)
  -http-port value
    	 (default 6000)
  -http-ssl

  -http-timeout value
    	 (default 15s)
exit status 2
```

## Options for flag tag

The flag default key string is the struct field name but can be specified in the struct field's tag value.
The "flag" key in the struct field's tag value is the key name, followed by an optional comma and options. Examples:
```golang
// Field is ignored by this package.
Field int `flag:"-"`

// Field appears in flags as "myName".
Field int `flag:"myName"`

// If this field is from nested struct, prefix from parent struct will be ingored.
Field int `flag:"~myName"`

// You can set short name for flags by providing it's value after a space
// Prefixes will not be applied for short names.
Field int `flag:"myName a"`

// this field will be removed from generated help text.
Field int `flag:",hidden"`

// this field will be marked as deprecated in generated help text
Field int `flag:",deprecated"`
```

## Options for desc tag
If you specify description in description tag (`desc` by default) it will be used in USAGE section.

```golang
Addr string `desc:"HTTP host"`
```sh
this description produces something like:
```
  -addr value
    	HTTP host (default 127.0.0.1)
```

## Options for env tag


## Options for Parse function:

```golang
// DescTag sets custom description tag. It is "desc" by default.
func DescTag(val string)

// FlagTag sets custom flag tag. It is "flag" be default.
func FlagTag(val string)

// Prefix sets prefix that will be applied for all flags (if they are not marked as ~).
func Prefix(val string)

// EnvPrefix sets prefix that will be applied for all environment variables (if they are not marked as ~).
func EnvPrefix(val string)

// FlagDivider sets custom divider for flags. It is dash by default. e.g. "flag-name".
func FlagDivider(val string)

// EnvDivider sets custom divider for environment variables.
// It is underscore by default. e.g. "ENV_NAME".
func EnvDivider(val string)

// Validator sets validator function for flags.
// Check existed validators in sflags/validator package.
func Validator(val ValidateFunc)

// Set to false if you don't want anonymous structure fields to be flatten.
func Flatten(val bool)
```


## Known issues

 - kingpin doesn't pass value for boolean arguments. Counter can't get initial value from arguments.
 
## Similar projects

 * https://github.com/jaffee/commandeer
 * https://github.com/anacrolix/tagflag
 * https://github.com/jessevdk/go-flags
