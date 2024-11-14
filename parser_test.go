package sflags

import (
	"errors"
	"net"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strP(value string) *string {
	return &value
}

type simple struct {
	Name string
}

func TestParseStruct(t *testing.T) {
	simpleCfg := &struct {
		Name  string `desc:"name description" env:"-"`
		Name2 string `flag:"name_two t,hidden,deprecated"`
		Name3 string `env:"NAME_THREE"`
		Name4 *string
		Name5 string `flag:"-"`
		name6 string
		Name7 int     `flag:",required"`
		Name8 float64 `flag:"name_8" env:"NAME_8,~nn_8"`

		Addr *net.TCPAddr

		Map map[string]int
	}{
		Name:  "name_value",
		Name2: "name2_value",
		Name4: strP("name_value4"),
		Addr: &net.TCPAddr{
			IP: net.ParseIP("127.0.0.1"),
		},
		name6: "name6_value",
		Map:   map[string]int{"test": 15},
	}
	diffTypesCfg := &struct {
		StringValue      string
		ByteValue        byte
		StringSliceValue []string
		BoolSliceValue   []bool
		CounterValue     Counter
		RegexpValue      *regexp.Regexp
		FuncValue        func() // will be ignored
		MapInt8Bool      map[int8]bool
		MapInt16Int8     map[int16]int8
		MapStringInt64   map[string]int64
		MapStringString  map[string]string
		MapBoolString    map[bool]string
	}{
		StringValue:      "string",
		ByteValue:        10,
		StringSliceValue: []string{},
		BoolSliceValue:   []bool{},
		CounterValue:     10,
		RegexpValue:      &regexp.Regexp{},
		MapStringInt64:   map[string]int64{"test": 888},
		MapStringString:  map[string]string{"test": "test-val"},
	}
	nestedCfg := &struct {
		Sub struct {
			Name  string `desc:"name description"`
			Name2 string `env:"NAME_TWO"`
			Name3 string `flag:"~name3" env:"~NAME_THREE"`
			SUB2  *struct {
				Name4 string
				Name5 string `env:"name_five"`
			}
		}
	}{
		Sub: struct {
			Name  string `desc:"name description"`
			Name2 string `env:"NAME_TWO"`
			Name3 string `flag:"~name3" env:"~NAME_THREE"`
			SUB2  *struct {
				Name4 string
				Name5 string `env:"name_five"`
			}
		}{
			Name:  "name_value",
			Name2: "name2_value",
			SUB2: &struct {
				Name4 string
				Name5 string `env:"name_five"`
			}{
				Name4: "name4_value",
			},
		},
	}
	descCfg := &struct {
		Name  string `desc:"name description"`
		Name2 string `description:"name2 description"`
	}{}
	anonymousCfg := &struct {
		Name1 string
		simple
	}{
		simple: simple{
			Name: "name_value",
		},
	}

	tt := []struct {
		name string

		cfg        interface{}
		optFuncs   []OptFunc
		expFlagSet []*Flag
		expErr     error
	}{
		{
			name: "SimpleCfg test",
			cfg:  simpleCfg,
			expFlagSet: []*Flag{
				{
					Name:     "name",
					EnvNames: nil,
					DefValue: "name_value",
					Value:    newStringValue(&simpleCfg.Name),
					Usage:    "name description",
				},
				{
					Name:       "name_two",
					Short:      "t",
					EnvNames:   []string{"NAME_TWO"},
					DefValue:   "name2_value",
					Value:      newStringValue(&simpleCfg.Name2),
					Hidden:     true,
					Deprecated: true,
				},
				{
					Name:     "name3",
					EnvNames: []string{"NAME_THREE"},
					DefValue: "",
					Value:    newStringValue(&simpleCfg.Name3),
				},
				{
					Name:     "name4",
					EnvNames: []string{"NAME4"},
					DefValue: "name_value4",
					Value:    newStringValue(simpleCfg.Name4),
				},
				{
					Name:     "name7",
					EnvNames: []string{"NAME7"},
					Required: true,
					DefValue: "0",
					Value:    newIntValue(&simpleCfg.Name7),
				},
				{
					Name:     "name_8",
					EnvNames: []string{"NAME_8", "nn_8"},
					DefValue: "0",
					Value:    newFloat64Value(&simpleCfg.Name8),
				},
				{
					Name:     "addr",
					EnvNames: []string{"ADDR"},
					DefValue: "127.0.0.1:0",
					Value:    newTCPAddrValue(simpleCfg.Addr),
				},
				{
					Name:     "map",
					EnvNames: []string{"MAP"},
					DefValue: "map[test:15]",
					Value:    newStringIntMapValue(&simpleCfg.Map),
				},
			},
		},
		{
			name:     "SimpleCfg test with custom env_prefix and divider",
			cfg:      simpleCfg,
			optFuncs: []OptFunc{EnvPrefix("PP|"), EnvDivider("|")},
			expFlagSet: []*Flag{
				{
					Name:     "name",
					EnvNames: nil,
					DefValue: "name_value",
					Value:    newStringValue(&simpleCfg.Name),
					Usage:    "name description",
				},
				{
					Name:       "name_two",
					Short:      "t",
					EnvNames:   []string{"PP|NAME_TWO"},
					DefValue:   "name2_value",
					Value:      newStringValue(&simpleCfg.Name2),
					Hidden:     true,
					Deprecated: true,
				},
				{
					Name:     "name3",
					EnvNames: []string{"PP|NAME_THREE"},
					DefValue: "",
					Value:    newStringValue(&simpleCfg.Name3),
				},
				{
					Name:     "name4",
					EnvNames: []string{"PP|NAME4"},
					DefValue: "name_value4",
					Value:    newStringValue(simpleCfg.Name4),
				},
				{
					Name:     "name7",
					EnvNames: []string{"PP|NAME7"},
					Required: true,
					DefValue: "0",
					Value:    newIntValue(&simpleCfg.Name7),
				},
				{
					Name:     "name_8",
					EnvNames: []string{"PP|NAME_8", "nn_8"},
					DefValue: "0",
					Value:    newFloat64Value(&simpleCfg.Name8),
				},
				{
					Name:     "addr",
					EnvNames: []string{"PP|ADDR"},
					DefValue: "127.0.0.1:0",
					Value:    newTCPAddrValue(simpleCfg.Addr),
				},
				{
					Name:     "map",
					EnvNames: []string{"PP|MAP"},
					DefValue: "map[test:15]",
					Value:    newStringIntMapValue(&simpleCfg.Map),
				},
			},
			expErr: nil,
		},
		{
			name: "DifferentTypesCfg",
			cfg:  diffTypesCfg,
			expFlagSet: []*Flag{
				{
					Name:     "string-value",
					EnvNames: []string{"STRING_VALUE"},
					DefValue: "string",
					Value:    newStringValue(&diffTypesCfg.StringValue),
					Usage:    "",
				},
				{
					Name:     "byte-value",
					EnvNames: []string{"BYTE_VALUE"},
					DefValue: "10",
					Value:    newUint8Value(&diffTypesCfg.ByteValue),
					Usage:    "",
				},
				{
					Name:     "string-slice-value",
					EnvNames: []string{"STRING_SLICE_VALUE"},
					DefValue: "[]",
					Value:    newStringSliceValue(&diffTypesCfg.StringSliceValue),
					Usage:    "",
				},
				{
					Name:     "bool-slice-value",
					EnvNames: []string{"BOOL_SLICE_VALUE"},
					DefValue: "[]",
					Value:    newBoolSliceValue(&diffTypesCfg.BoolSliceValue),
					Usage:    "",
				},
				{
					Name:     "counter-value",
					EnvNames: []string{"COUNTER_VALUE"},
					DefValue: "10",
					Value:    &diffTypesCfg.CounterValue,
					Usage:    "",
				},
				{
					Name:     "regexp-value",
					EnvNames: []string{"REGEXP_VALUE"},
					DefValue: "",
					Value:    newRegexpValue(&diffTypesCfg.RegexpValue),
					Usage:    "",
				},
				{
					Name:     "map-int8-bool",
					EnvNames: []string{"MAP_INT8_BOOL"},
					DefValue: "",
					Value:    newInt8BoolMapValue(&diffTypesCfg.MapInt8Bool),
				},
				{
					Name:     "map-int16-int8",
					EnvNames: []string{"MAP_INT16_INT8"},
					DefValue: "",
					Value:    newInt16Int8MapValue(&diffTypesCfg.MapInt16Int8),
				},
				{
					Name:     "map-string-int64",
					EnvNames: []string{"MAP_STRING_INT64"},
					DefValue: "map[test:888]",
					Value:    newStringInt64MapValue(&diffTypesCfg.MapStringInt64),
				},
				{
					Name:     "map-string-string",
					EnvNames: []string{"MAP_STRING_STRING"},
					DefValue: "map[test:test-val]",
					Value:    newStringStringMapValue(&diffTypesCfg.MapStringString),
				},
			},
		},
		{
			name: "NestedCfg",
			cfg:  nestedCfg,
			expFlagSet: []*Flag{
				{
					Name:     "sub-name",
					EnvNames: []string{"SUB_NAME"},
					DefValue: "name_value",
					Value:    newStringValue(&nestedCfg.Sub.Name),
					Usage:    "name description",
				},
				{
					Name:     "sub-name2",
					EnvNames: []string{"SUB_NAME_TWO"},
					DefValue: "name2_value",
					Value:    newStringValue(&nestedCfg.Sub.Name2),
				},
				{
					Name:     "name3",
					EnvNames: []string{"NAME_THREE"},
					DefValue: "",
					Value:    newStringValue(&nestedCfg.Sub.Name3),
				},
				{
					Name:     "sub-sub2-name4",
					EnvNames: []string{"SUB_SUB2_NAME4"},
					DefValue: "name4_value",
					Value:    newStringValue(&nestedCfg.Sub.SUB2.Name4),
				},
				{
					Name:     "sub-sub2-name5",
					EnvNames: []string{"SUB_SUB2_name_five"},
					DefValue: "",
					Value:    newStringValue(&nestedCfg.Sub.SUB2.Name5),
				},
			},
			expErr: nil,
		},
		{
			name:     "DescCfg with custom desc tag",
			cfg:      descCfg,
			optFuncs: []OptFunc{DescTag("description")},
			expFlagSet: []*Flag{
				{
					Name:     "name",
					EnvNames: []string{"NAME"},
					Value:    newStringValue(&descCfg.Name),
				},
				{
					Name:     "name2",
					EnvNames: []string{"NAME2"},
					Value:    newStringValue(&descCfg.Name2),
					Usage:    "name2 description",
				},
			},
		},
		{
			name: "Anonymoust cfg with disabled flatten",
			cfg:  anonymousCfg,
			expFlagSet: []*Flag{
				{
					Name:     "name1",
					EnvNames: []string{"NAME1"},
					Value:    newStringValue(&anonymousCfg.Name1),
				},
				{
					Name:     "name",
					EnvNames: []string{"NAME"},
					DefValue: "name_value",
					Value:    newStringValue(&anonymousCfg.Name),
				},
			},
		},
		{
			name:     "Anonymoust cfg with enabled flatten",
			cfg:      anonymousCfg,
			optFuncs: []OptFunc{Flatten(false)},
			expFlagSet: []*Flag{
				{
					Name:     "name1",
					EnvNames: []string{"NAME1"},
					Value:    newStringValue(&anonymousCfg.Name1),
				},
				{
					Name:     "simple-name",
					EnvNames: []string{"SIMPLE_NAME"},
					DefValue: "name_value",
					Value:    newStringValue(&anonymousCfg.Name),
				},
			},
		},
		{
			name: "We need pointer to structure",
			cfg: struct {
			}{},
			expErr: errors.New("object must be a pointer to struct or interface"),
		},
		{
			name:   "We need pointer to structure 2",
			cfg:    strP("something"),
			expErr: errors.New("object must be a pointer to struct or interface"),
		},
		{
			name:   "We need non nil object",
			cfg:    nil,
			expErr: errors.New("object cannot be nil"),
		},
		{
			name:   "We need non nil value",
			cfg:    (*simple)(nil),
			expErr: errors.New("object cannot be nil"),
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			flagSet, err := ParseStruct(test.cfg, test.optFuncs...)
			if test.expErr == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, test.expErr, err)
			}
			assert.Equal(t, test.expFlagSet, flagSet)
		})
	}
}

func TestParseStruct_NilValue(t *testing.T) {
	name2Value := "name2_value"
	cfg := struct {
		Name1  *string
		Name2  *string
		Regexp *regexp.Regexp
	}{
		Name2: &name2Value,
	}
	assert.Nil(t, cfg.Name1)
	assert.Nil(t, cfg.Regexp)
	assert.NotNil(t, cfg.Name2)

	flags, err := ParseStruct(&cfg)
	require.NoError(t, err)
	require.Equal(t, 3, len(flags))
	assert.NotNil(t, cfg.Name1)
	assert.NotNil(t, cfg.Name2)
	assert.NotNil(t, cfg.Regexp)
	assert.Equal(t, name2Value, flags[1].Value.(Getter).Get())

	err = flags[0].Value.Set("name1value")
	require.NoError(t, err)
	assert.Equal(t, "name1value", *cfg.Name1)

	err = flags[2].Value.Set("aabbcc")
	require.NoError(t, err)
	assert.Equal(t, "aabbcc", cfg.Regexp.String())
}

func TestParseStruct_WithValidator(t *testing.T) {
	var cfg simple

	testErr := errors.New("validator test error")

	validator := Validator(func(val string, field reflect.StructField, cfg interface{}) error {
		return testErr
	})

	flags, err := ParseStruct(&cfg, validator)
	require.NoError(t, err)
	require.Equal(t, 1, len(flags))
	assert.NotNil(t, cfg.Name)

	err = flags[0].Value.Set("aabbcc")
	require.Error(t, err)
	assert.Equal(t, testErr, err)
}

func TestFlagDivider(t *testing.T) {
	opt := opts{
		flagDivider: "-",
	}
	FlagDivider("_")(&opt)
	assert.Equal(t, "_", opt.flagDivider)
}

func TestFlagTag(t *testing.T) {
	opt := opts{
		flagTag: "flags",
	}
	FlagTag("superflag")(&opt)
	assert.Equal(t, "superflag", opt.flagTag)
}

func TestValidator(t *testing.T) {
	opt := opts{
		validator: nil,
	}
	Validator(func(string, reflect.StructField, interface{}) error {
		return nil
	})(&opt)
	assert.NotNil(t, opt.validator)
}

func TestFlatten(t *testing.T) {
	opt := opts{
		flatten: true,
	}
	Flatten(false)(&opt)
	assert.Equal(t, false, opt.flatten)
}
