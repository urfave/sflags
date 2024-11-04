package gcli

import (
	"errors"
	"io"
	"testing"

	"github.com/octago/sflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

type cfg1 struct {
	StringValue1 string
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
		expErr2 error // cli Parse error
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
			args: []string{},
		},
		{
			name: "Test cfg1 short option",
			cfg: &cfg1{
				StringValue2: "string_value2_value",
			},
			expCfg: &cfg1{
				StringValue2: "string_value2_value2",
			},
			args: []string{
				"-s=string_value2_value2",
			},
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
				"--string-value1", "string_value1_value2",
				"--string-value-two", "string_value2_value2",
				"--counter-value1=2", "--counter-value1",
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
			expErr2: errors.New("flag provided but not defined: -bad-value"),
		},
		{
			name:    "Test bad cfg value",
			cfg:     "bad config",
			expErr1: errors.New("object must be a pointer to struct or interface"),
		},
	}
	// forbid urfave/cli to exit
	cli.OsExiter = func(i int) {}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags, err := Parse(test.cfg)
			if test.expErr1 != nil {
				require.Error(t, err)
				require.Equal(t, test.expErr1, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}
			cliApp := cli.NewApp()
			cliApp.Action = func(c *cli.Context) error {
				return nil
			}
			cliApp.UseShortOptionHandling = true
			cli.ErrWriter = io.Discard
			cliApp.OnUsageError = func(_ *cli.Context, err error, _ bool) error {
				return err
			}

			cliApp.Flags = flags
			args := append([]string{"cliApp"}, test.args...)
			err = cliApp.Run(args)
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
