package gcli

import (
	"github.com/octago/sflags"
	"github.com/urfave/cli/v2"
)

// GenerateTo takes a list of sflag.Flag,
// that are parsed from some config structure, and put it to dst.
func GenerateTo(src []*sflags.Flag, dst *[]cli.Flag) {
	for _, srcFlag := range src {
		name := srcFlag.Name
		if srcFlag.Short != "" {
			name += ", " + srcFlag.Short
		}
		*dst = append(*dst, &cli.GenericFlag{
			Name:   name,
			EnvVar: srcFlag.EnvName,
			Hidden: srcFlag.Hidden,
			Usage:  srcFlag.Usage,
			Value:  srcFlag.Value,
		})
	}
}

// ParseTo parses cfg, that is a pointer to some structure,
// and puts it to dst.
func ParseTo(cfg interface{}, dst *[]cli.Flag, optFuncs ...sflags.OptFunc) error {
	flags, err := sflags.ParseStruct(cfg, optFuncs...)
	if err != nil {
		return err
	}
	GenerateTo(flags, dst)
	return nil
}

// Parse parses cfg, that is a pointer to some structure,
// puts it to the new flag.FlagSet and returns it.
func Parse(cfg interface{}, optFuncs ...sflags.OptFunc) ([]cli.Flag, error) {
	flags := make([]cli.Flag, 0)
	err := ParseTo(cfg, &flags, optFuncs...)
	if err != nil {
		return nil, err
	}
	return flags, nil
}
