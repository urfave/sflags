package gkingpin

import (
	"errors"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/sflags"
)

type cfg1 struct {
	StringValue1 string `flag:",required"`
	StringValue2 string `flag:"string-value-two s"`

	CounterValue1 sflags.Counter

	StringSliceValue1 []string
}

func TestParse(t *testing.T) {
	tests := []struct {
		name string

		cfg     interface{}
		args    []string
		expCfg  interface{}
		expErr1 error // sflag Parse error
		expErr2 error // kingpin Parse error
	}{
		{
			name: "Test cfg1",
			cfg: &cfg1{
				StringValue1: "string_value1_value",
				StringValue2: "string_value2_value",

				CounterValue1: 1,

				StringSliceValue1: []string{"one", "two"},
			},
			expCfg: &cfg1{
				StringValue1: "string_value1_value2",
				StringValue2: "string_value2_value2",

				CounterValue1: 3,

				StringSliceValue1: []string{
					"one2", "two2", "three", "4"},
			},
			args: []string{
				"--string-value1", "string_value1_value2",
				"--string-value-two", "string_value2_value2",
				"--counter-value1", "--counter-value1",
				"--string-slice-value1", "one2",
				"--string-slice-value1", "two2",
				"--string-slice-value1", "three,4",
			},
		},
		{
			name: "Test cfg1 no args",
			cfg: &cfg1{
				StringValue1: "string_value1_value",
				StringValue2: "",
			},
			expCfg: &cfg1{
				StringValue1: "string_value1_value",
				StringValue2: "",
			},
			args:    []string{},
			expErr2: errors.New("required flag(s) '--string-value1' not provided"),
		},
		{
			name: "Test cfg1 short option",
			cfg: &cfg1{
				StringValue2: "string_value2_value",
			},
			expCfg: &cfg1{
				StringValue1: "string_value1_value2",
				StringValue2: "string_value2_value2",
			},
			args: []string{
				"-s", "string_value2_value2",
				"--string-value1", "string_value1_value2",
			},
		},
		{
			name: "Test cfg1 without default values",
			cfg:  &cfg1{},
			expCfg: &cfg1{
				StringValue1: "string_value1_value2",
				StringValue2: "string_value2_value2",

				CounterValue1: 1,
			},
			args: []string{
				"--string-value1", "string_value1_value2",
				"--string-value-two", "string_value2_value2",
				// kingpin can't pass value for boolean arguments.
				//"--counter-value1", "2",
				"--counter-value1",
			},
		},
		{
			name: "Test cfg1 bad option",
			cfg: &cfg1{
				StringValue1: "string_value1_value",
			},
			args: []string{
				"--bad-value=string_value1_value2",
			},
			expErr2: errors.New("unknown long flag '--bad-value'"),
		},
		{
			name:    "Test bad cfg value",
			cfg:     "bad config",
			expErr1: errors.New("object must be a pointer to struct or interface"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := kingpin.New("testApp", "")
			app.Terminate(nil)

			err := ParseTo(test.cfg, app)
			if test.expErr1 != nil {
				require.Error(t, err)
				require.Equal(t, test.expErr1, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}

			_, err = app.Parse(test.args)
			if test.expErr2 != nil {
				require.Error(t, err)
				require.Equal(t, test.expErr2, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}
			assert.Equal(t, test.expCfg, test.cfg)
		})
	}
}
