package gflag

import (
	"errors"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/sflags"
)

type cfg1 struct {
	StringValue1 string
	StringValue2 string `flag:"string-value-two"`

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
		expErr2 error // flag Parse error
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
				"-string-value1", "string_value1_value2",
				"-string-value-two", "string_value2_value2",
				"-counter-value1", "-counter-value1",
				"-string-slice-value1", "one2",
				"-string-slice-value1", "two2",
				"-string-slice-value1", "three,4",
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
			args: []string{},
		},
		{
			name: "Test cfg1 without default values",
			cfg:  &cfg1{},
			expCfg: &cfg1{
				StringValue1: "string_value1_value2",
				StringValue2: "string_value2_value2",

				CounterValue1: 3,
			},
			args: []string{
				"-string-value1", "string_value1_value2",
				"-string-value-two", "string_value2_value2",
				"-counter-value1=2", "-counter-value1",
			},
		},
		{
			name:    "Test bad cfg value",
			cfg:     "bad config",
			expErr1: errors.New("object must be a pointer to struct or interface"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs, err := Parse(test.cfg)
			require.Equal(t, test.expErr1, err)
			if err != nil {
				return
			}
			err = fs.Parse(test.args)
			require.Equal(t, test.expErr2, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.expCfg, test.cfg)
		})
	}
}

func TestParseToDef(t *testing.T) {
	oldCommandLine := flag.CommandLine
	defer func() {
		flag.CommandLine = oldCommandLine
	}()
	cfg := &cfg1{StringValue1: "value1"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	err := ParseToDef(cfg)
	assert.NoError(t, err)
	err = flag.CommandLine.Parse([]string{"-string-value1", "value2"})
	assert.NoError(t, err)
	assert.Equal(t, "value2", cfg.StringValue1)
	err = ParseToDef("bad string")
	assert.Error(t, err)
}
