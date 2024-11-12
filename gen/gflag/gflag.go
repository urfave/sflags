package gflag

import (
	"flag"
	"os"

	"github.com/urfave/sflags"
)

// flagSet describes interface,
// that's implemented by flag library and required by sflags.
type flagSet interface {
	Var(value flag.Value, name string, usage string)
}

var _ flagSet = (*flag.FlagSet)(nil)

// GenerateTo takes a list of sflag.Flag,
// that are parsed from some config structure, and put it to dst.
func GenerateTo(src []*sflags.Flag, dst flagSet) {
	for _, srcFlag := range src {
		dst.Var(srcFlag.Value, srcFlag.Name, srcFlag.Usage)
	}
}

// ParseTo parses cfg, that is a pointer to some structure,
// and puts it to dst.
func ParseTo(cfg interface{}, dst flagSet, optFuncs ...sflags.OptFunc) error {
	flags, err := sflags.ParseStruct(cfg, optFuncs...)
	if err != nil {
		return err
	}
	GenerateTo(flags, dst)
	return nil
}

// Parse parses cfg, that is a pointer to some structure,
// puts it to the new flag.FlagSet and returns it.
func Parse(cfg interface{}, optFuncs ...sflags.OptFunc) (*flag.FlagSet, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	err := ParseTo(cfg, fs, optFuncs...)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

// ParseToDef parses cfg, that is a pointer to some structure and
// puts it to the default flag.CommandLine.
func ParseToDef(cfg interface{}, optFuncs ...sflags.OptFunc) error {
	err := ParseTo(cfg, flag.CommandLine, optFuncs...)
	if err != nil {
		return err
	}
	return nil
}
