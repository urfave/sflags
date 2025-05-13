package sflags

import (
	"errors"
	"reflect"
	"strings"
)

const (
	defaultDescTag           = "desc"
	defaultFlagTag           = "flag"
	defaultEnvTag            = "env"
	defaultFlagDivider       = "-"
	defaultEnvDivider        = "_"
	defaultFlatten           = true
	defaultInheritHidden     = false
	defaultHidden            = false
	defaultInheritDeprecated = false
	defaultDeprecated        = false
)

// ValidateFunc describes a validation func,
// that takes string val for flag from command line,
// field that's associated with this flag in structure cfg.
// Should return error if validation fails.
type ValidateFunc func(val string, field reflect.StructField, cfg interface{}) error

type opts struct {
	descTag           string
	flagTag           string
	prefix            string
	envPrefix         string
	flagDivider       string
	envDivider        string
	flatten           bool
	validator         ValidateFunc
	inheritHidden     bool
	hidden            bool
	inheritDeprecated bool
	deprecated        bool
}

func (o opts) apply(optFuncs ...OptFunc) opts {
	for _, optFunc := range optFuncs {
		optFunc(&o)
	}
	return o
}

// OptFunc sets values in opts structure.
type OptFunc func(opt *opts)

// DescTag sets custom description tag. It is "desc" by default.
func DescTag(val string) OptFunc { return func(opt *opts) { opt.descTag = val } }

// FlagTag sets custom flag tag. It is "flag" be default.
func FlagTag(val string) OptFunc { return func(opt *opts) { opt.flagTag = val } }

// Prefix sets prefix that will be applied for all flags (if they are not marked as ~).
func Prefix(val string) OptFunc { return func(opt *opts) { opt.prefix = val } }

// InheritHidden enables inheriting the hidden flag for all nested flags if set for a parent flag
func InheritHidden() OptFunc { return func(opt *opts) { opt.inheritHidden = true } }

// hidden sets the hidden flag for all nested flags if set for a parent flag
func hidden(val bool) OptFunc {
	return func(opt *opts) { opt.hidden = val }
}

// InheritDeprecated enables inheriting the deprecated flag for all nested flags if set for a parent flag
func InheritDeprecated() OptFunc { return func(opt *opts) { opt.inheritDeprecated = true } }

// deprecated sets the deprecated flag for all nested flags if set for a parent flag
func deprecated(val bool) OptFunc {
	return func(opt *opts) { opt.deprecated = val }
}

// EnvPrefix sets prefix that will be applied for all environment variables (if they are not marked as ~).
func EnvPrefix(val string) OptFunc { return func(opt *opts) { opt.envPrefix = val } }

// FlagDivider sets custom divider for flags. It is dash by default. e.g. "flag-name".
func FlagDivider(val string) OptFunc { return func(opt *opts) { opt.flagDivider = val } }

// EnvDivider sets custom divider for environment variables.
// It is underscore by default. e.g. "ENV_NAME".
func EnvDivider(val string) OptFunc { return func(opt *opts) { opt.envDivider = val } }

// Validator sets validator function for flags.
// Check existed validators in sflags/validator package.
func Validator(val ValidateFunc) OptFunc { return func(opt *opts) { opt.validator = val } }

// Flatten set flatten option.
// Set to false if you don't want anonymous structure fields to be flatten.
func Flatten(val bool) OptFunc { return func(opt *opts) { opt.flatten = val } }

func copyOpts(val opts) OptFunc { return func(opt *opts) { *opt = val } }

func hasOption(options []string, option string) bool {
	for _, opt := range options {
		if opt == option {
			return true
		}
	}
	return false
}

func defOpts() opts {
	return opts{
		descTag:           defaultDescTag,
		flagTag:           defaultFlagTag,
		flagDivider:       defaultFlagDivider,
		envDivider:        defaultEnvDivider,
		flatten:           defaultFlatten,
		inheritHidden:     defaultInheritHidden,
		hidden:            defaultHidden,
		inheritDeprecated: defaultInheritDeprecated,
		deprecated:        defaultDeprecated,
	}
}

func parseFlagTag(field reflect.StructField, opt opts) *Flag {
	flag := Flag{}
	ignoreFlagPrefix := false
	flag.Name = camelToFlag(field.Name, opt.flagDivider)
	if flagTags := strings.Split(field.Tag.Get(opt.flagTag), ","); len(flagTags) > 0 {
		switch fName := flagTags[0]; fName {
		case "-":
			return nil
		case "":
		default:
			fNameSplitted := strings.Split(fName, " ")
			if len(fNameSplitted) > 1 {
				fName = fNameSplitted[0]
				flag.Short = fNameSplitted[1]
			}
			if strings.HasPrefix(fName, "~") {
				flag.Name = fName[1:]
				ignoreFlagPrefix = true
			} else {
				flag.Name = fName
			}
		}
		flag.Hidden = hasOption(flagTags[1:], "hidden")
		flag.Deprecated = hasOption(flagTags[1:], "deprecated")
		flag.Required = hasOption(flagTags[1:], "required")
	}

	if opt.prefix != "" && !ignoreFlagPrefix {
		flag.Name = opt.prefix + flag.Name
	}

	if opt.deprecated {
		flag.Deprecated = opt.deprecated
	}
	if opt.hidden {
		flag.Hidden = opt.hidden
	}
	return &flag
}

func parseEnv(flagName string, field reflect.StructField, opt opts) []string {
	var envVars []string
	flagEnvVar := flagToEnv(flagName, opt.flagDivider, opt.envDivider)
	if envTags := strings.Split(field.Tag.Get(defaultEnvTag), ","); len(envTags) > 0 {
		switch envName := envTags[0]; envName {
		case "-":
			return envVars
		case "":
			// if tag is `env:""` then env var will be taken from flag name
			envVars = append(envVars, opt.envPrefix+flagEnvVar)
		default:
			// if tag is `env:"NAME"` then env var is envPrefix_flagPrefix_NAME
			// if tag is `env:"~NAME"` then env var is NAME
			for _, envName := range envTags {
				ignoreEnvPrefix := false
				var envVar string
				if strings.HasPrefix(envName, "~") {
					envVar = envName[1:]
					ignoreEnvPrefix = true
				} else {
					envVar = envName
					if opt.prefix != "" {
						envVar = flagToEnv(
							opt.prefix,
							opt.flagDivider,
							opt.envDivider) + envVar
					}
				}
				if envVar != "" {
					if !ignoreEnvPrefix {
						envVars = append(envVars, opt.envPrefix+envVar)
					} else {
						envVars = append(envVars, envVar)
					}
				}
			}
		}
	} else if flagEnvVar != "" {
		envVars = append(envVars, opt.envPrefix+flagEnvVar)
	}
	return envVars
}

// ParseStruct parses structure and returns list of flags based on this structure.
// This list of flags can be used by generators for flag, kingpin, cobra, pflag, urfave/cli.
func ParseStruct(cfg interface{}, optFuncs ...OptFunc) ([]*Flag, error) {
	// what we want is Ptr to Structure
	if cfg == nil {
		return nil, errors.New("object cannot be nil")
	}
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("object must be a pointer to struct or interface")
	}
	if v.IsNil() {
		return nil, errors.New("object cannot be nil")
	}
	switch e := v.Elem(); e.Kind() {
	case reflect.Struct:
		return parseStruct(e, optFuncs...), nil
	default:
		return nil, errors.New("object must be a pointer to struct or interface")
	}
}

func parseVal(value reflect.Value, optFuncs ...OptFunc) ([]*Flag, Value) {
	// value is addressable, let's check if we can parse it
	if value.CanAddr() && value.Addr().CanInterface() {
		valueInterface := value.Addr().Interface()
		val := parseGenerated(valueInterface)
		if val != nil {
			return nil, val
		}
		// check if field implements Value interface
		if val, casted := valueInterface.(Value); casted {
			return nil, val
		}
	}

	switch value.Kind() {
	case reflect.Ptr:
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		val := parseGeneratedPtrs(value.Addr().Interface())
		if val != nil {
			return nil, val
		}
		return parseVal(value.Elem(), optFuncs...)
	case reflect.Struct:
		flags := parseStruct(value, optFuncs...)
		return flags, nil
	case reflect.Map:
		mapType := value.Type()
		keyKind := value.Type().Key().Kind()

		// check that map key is string or integer
		if !anyOf(MapAllowedKinds, keyKind) {
			break
		}

		if value.IsNil() {
			value.Set(reflect.MakeMap(mapType))
		}

		valueInterface := value.Addr().Interface()
		val := parseGeneratedMap(valueInterface)
		if val != nil {
			return nil, val
		}
	}
	return nil, nil
}

func parseStruct(value reflect.Value, optFuncs ...OptFunc) []*Flag {
	opt := defOpts().apply(optFuncs...)

	flags := []*Flag{}

	valueType := value.Type()
fields:
	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		fieldValue := value.Field(i)
		// skip unexported and non anonymous fields
		if field.PkgPath != "" && !field.Anonymous {
			continue fields
		}

		flag := parseFlagTag(field, opt)
		if flag == nil {
			continue fields
		}

		flag.EnvNames = parseEnv(flag.Name, field, opt)
		flag.Usage = field.Tag.Get(opt.descTag)
		prefix := flag.Name + opt.flagDivider
		if field.Anonymous && opt.flatten {
			prefix = opt.prefix
		}

		nestedOpts := []OptFunc{copyOpts(opt), Prefix(prefix)}
		if opt.inheritHidden {
			nestedOpts = append(nestedOpts, hidden(flag.Hidden))
		}
		if opt.inheritDeprecated {
			nestedOpts = append(nestedOpts, deprecated(flag.Deprecated))
		}

		nestedFlags, val := parseVal(fieldValue,
			nestedOpts...,
		)

		// field contains a simple value.
		if val != nil {
			if opt.validator != nil {
				val = &validateValue{
					Value: val,
					validateFunc: func(val string) error {
						return opt.validator(val, field, value.Interface())
					},
				}
			}
			flag.Value = val
			flag.DefValue = val.String()
			flags = append(flags, flag)
			continue fields
		}
		// field is a structure
		if len(nestedFlags) > 0 {
			flags = append(flags, nestedFlags...)
			continue fields
		}

	}
	return flags
}

func anyOf(kinds []reflect.Kind, needle reflect.Kind) bool {
	for _, kind := range kinds {
		if kind == needle {
			return true
		}
	}

	return false
}
