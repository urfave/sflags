package gcli

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/octago/sflags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

type cfg2 struct {
	StringValue1 string
	StringValue2 string `flag:"string-value-two s"`

	CounterValue1 sflags.Counter

	StringSliceValue1 []string
}

func TestParseV3(t *testing.T) {
	tests := []struct {
		name string

		cfg     interface{}
		args    []string
		expCfg  interface{}
		expErr1 error // sflag Parse error
		expErr2 error // cli Parse error
	}{
		{
			name: "Test cfg2",
			cfg: &cfg2{
				StringValue1: "string_value1_value",
				StringValue2: "string_value2_value",

				CounterValue1: 1,

				StringSliceValue1: []string{"one", "two"},
			},
			expCfg: &cfg2{
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
			name: "Test cfg2 no args",
			cfg: &cfg2{
				StringValue1: "string_value1_value",
				StringValue2: "",
			},
			expCfg: &cfg2{
				StringValue1: "string_value1_value",
				StringValue2: "",
			},
			args: []string{},
		},
		{
			name: "Test cfg2 short option",
			cfg: &cfg2{
				StringValue2: "string_value2_value",
			},
			expCfg: &cfg2{
				StringValue2: "string_value2_value2",
			},
			args: []string{
				"-s=string_value2_value2",
			},
		},
		{
			name: "Test cfg2 without default values",
			cfg:  &cfg2{},
			expCfg: &cfg2{
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
			name: "Test cfg2 bad option",
			cfg: &cfg2{
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
			flags, err := ParseV3(test.cfg)
			if test.expErr1 != nil {
				require.Error(t, err)
				require.Equal(t, test.expErr1, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}
			cmd := &cli.Command{}
			cmd.Action = func(_ context.Context, c *cli.Command) error {
				return nil
			}
			cmd.UseShortOptionHandling = true
			cli.ErrWriter = io.Discard
			cmd.OnUsageError = func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
				return err
			}

			cmd.Flags = flags
			args := append([]string{"cliApp"}, test.args...)
			err = cmd.Run(context.Background(), args)
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
