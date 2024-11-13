package gkingpin

import (
	"unicode/utf8"

	"github.com/alecthomas/kingpin/v2"
	"github.com/urfave/sflags"
)

type flagger interface {
	Flag(name, help string) *kingpin.FlagClause
}

// GenerateTo takes a list of sflag.Flag,
// that are parsed from some config structure, and put it to dst.
func GenerateTo(src []*sflags.Flag, dst flagger) {
	for _, srcFlag := range src {
		flag := dst.Flag(srcFlag.Name, srcFlag.Usage)
		flag.SetValue(srcFlag.Value)
		if srcFlag.EnvName != "" {
			flag.Envar(srcFlag.EnvName)
		}
		if srcFlag.Hidden {
			flag.Hidden()
		}
		if srcFlag.Required {
			flag.Required()
		}
		if srcFlag.Short != "" {
			r, _ := utf8.DecodeRuneInString(srcFlag.Short)
			if r != utf8.RuneError {
				flag.Short(r)
			}
		}

	}
}

// ParseTo parses cfg, that is a pointer to some structure,
// and puts it to dst.
func ParseTo(cfg interface{}, dst flagger, optFuncs ...sflags.OptFunc) error {
	flags, err := sflags.ParseStruct(cfg, optFuncs...)
	if err != nil {
		return err
	}
	GenerateTo(flags, dst)
	return nil
}
