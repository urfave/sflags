package gcli

import (
	"github.com/octago/sflags"
	"github.com/urfave/cli/v3"
)

type boolFlag interface {
	IsBoolFlag() bool
}

type value struct {
	v sflags.Value
}

func (v value) Get() any {
	return v.v
}

func (v value) Set(s string) error {
	return v.v.Set(s)
}

func (v value) String() string {
	return v.v.String()
}

func (v value) IsBoolFlag() bool {
	b, ok := v.v.(boolFlag)
	return ok && b.IsBoolFlag()
}

// GenerateToV3 takes a list of sflag.Flag,
// that are parsed from some config structure, and put it to dst.
func GenerateToV3(src []*sflags.Flag, dst *[]cli.Flag) {
	for _, srcFlag := range src {
		name := srcFlag.Name
		var aliases []string
		if srcFlag.Short != "" {
			aliases = append(aliases, srcFlag.Short)
		}
		*dst = append(*dst, &cli.GenericFlag{
			Name:    name,
			Sources: cli.EnvVars(srcFlag.EnvName),
			Aliases: aliases,
			Hidden:  srcFlag.Hidden,
			Usage:   srcFlag.Usage,
			Value: &value{
				v: srcFlag.Value,
			},
		})
	}
}

// ParseToV3 parses cfg, that is a pointer to some structure,
// and puts it to dst.
func ParseToV3(cfg interface{}, dst *[]cli.Flag, optFuncs ...sflags.OptFunc) error {
	flags, err := sflags.ParseStruct(cfg, optFuncs...)
	if err != nil {
		return err
	}
	GenerateToV3(flags, dst)
	return nil
}

// ParseV3 parses cfg, that is a pointer to some structure,
// puts it to the new flag.FlagSet and returns it.
func ParseV3(cfg interface{}, optFuncs ...sflags.OptFunc) ([]cli.Flag, error) {
	flags := make([]cli.Flag, 0)
	err := ParseToV3(cfg, &flags, optFuncs...)
	if err != nil {
		return nil, err
	}
	return flags, nil
}
