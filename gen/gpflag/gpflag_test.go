package gpflag

import (
	"errors"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/sflags"
)

type cfg1 struct {
	StringValue1 string
	StringValue2 string `flag:"string-value-two s"`

	CounterValue1 sflags.Counter

	StringSliceValue1 []string
	DeprecatedValue1  string `flag:",deprecated" desc:"DEP_MESSAGE"`
}

func TestParse(t *testing.T) {
	tests := []struct {
		name string

		cfg     interface{}
		args    []string
		expCfg  interface{}
		expErr1 error // sflag Parse error
		expErr2 error // pflag Parse error
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
			expErr2: errors.New("unknown flag: --bad-value"),
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
			if test.expErr1 != nil {
				require.Error(t, err)
				require.Equal(t, test.expErr1, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}
			fs.Init("pflagTest", pflag.ContinueOnError)
			fs.SetOutput(io.Discard)
			err = fs.Parse(test.args)
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

func TestParseToDef(t *testing.T) {
	oldCommandLine := pflag.CommandLine
	defer func() {
		pflag.CommandLine = oldCommandLine
	}()
	cfg := &cfg1{StringValue1: "value1"}
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	err := ParseToDef(cfg)
	assert.NoError(t, err)
	err = pflag.CommandLine.Parse([]string{"--string-value1", "value2"})
	assert.NoError(t, err)
	assert.Equal(t, "value2", cfg.StringValue1)
	err = ParseToDef("bad string")
	assert.Error(t, err)
}

func TestPFlagGetters(t *testing.T) {
	// Test that pflag getter functions like GetInt work as expected.
	_, ipNet, err := net.ParseCIDR("127.0.0.1/24")
	require.NoError(t, err)

	cfg := &struct {
		IntValue   int
		Int8Value  int8
		Int32Value int32
		Int64Value int64

		UintValue   uint
		Uint8Value  uint8
		Uint16Value uint16
		Uint32Value uint32
		Uint64Value uint64

		Float32Value float32
		Float64Value float64

		BoolValue     bool
		StringValue   string
		DurationValue time.Duration
		CountValue    sflags.Counter

		IPValue    net.IP
		IPNetValue net.IPNet

		StringSliceValue []string
		IntSliceValue    []int
	}{
		IntValue:   10,
		Int8Value:  11,
		Int32Value: 12,
		Int64Value: 13,

		UintValue:   14,
		Uint8Value:  15,
		Uint16Value: 16,
		Uint32Value: 17,
		Uint64Value: 18,

		Float32Value: 19.1,
		Float64Value: 20.1,

		BoolValue:     true,
		StringValue:   "stringValue",
		DurationValue: time.Second * 10,
		CountValue:    30,

		IPValue:    net.ParseIP("127.0.0.1"),
		IPNetValue: *ipNet,

		StringSliceValue: []string{"one", "two"},
		IntSliceValue:    []int{10, 20},
	}
	flagSet, err := Parse(cfg)
	require.NoError(t, err)

	intValue, err := flagSet.GetInt("int-value")
	assert.NoError(t, err)
	assert.Equal(t, 10, intValue)

	int8Value, err := flagSet.GetInt8("int8-value")
	assert.NoError(t, err)
	assert.Equal(t, int8(11), int8Value)

	int32Value, err := flagSet.GetInt32("int32-value")
	assert.NoError(t, err)
	assert.Equal(t, int32(12), int32Value)

	int64Value, err := flagSet.GetInt64("int64-value")
	assert.NoError(t, err)
	assert.Equal(t, int64(13), int64Value)

	uintValue, err := flagSet.GetUint("uint-value")
	assert.NoError(t, err)
	assert.Equal(t, uint(14), uintValue)

	uint8Value, err := flagSet.GetUint8("uint8-value")
	assert.NoError(t, err)
	assert.Equal(t, uint8(15), uint8Value)

	uint16Value, err := flagSet.GetUint16("uint16-value")
	assert.NoError(t, err)
	assert.Equal(t, uint16(16), uint16Value)

	uint32Value, err := flagSet.GetUint32("uint32-value")
	assert.NoError(t, err)
	assert.Equal(t, uint32(17), uint32Value)

	uint64Value, err := flagSet.GetUint64("uint64-value")
	assert.NoError(t, err)
	assert.Equal(t, uint64(18), uint64Value)

	float32Value, err := flagSet.GetFloat32("float32-value")
	assert.NoError(t, err)
	assert.Equal(t, float32(19.1), float32Value)

	float64Value, err := flagSet.GetFloat64("float64-value")
	assert.NoError(t, err)
	assert.Equal(t, float64(20.1), float64Value)

	boolValue, err := flagSet.GetBool("bool-value")
	assert.NoError(t, err)
	assert.Equal(t, true, boolValue)

	countValue, err := flagSet.GetCount("count-value")
	assert.NoError(t, err)
	assert.Equal(t, 30, countValue)

	durationValue, err := flagSet.GetDuration("duration-value")
	assert.NoError(t, err)
	assert.Equal(t, time.Second*10, durationValue)

	stringValue, err := flagSet.GetString("string-value")
	assert.NoError(t, err)
	assert.Equal(t, "stringValue", stringValue)

	ipValue, err := flagSet.GetIP("ip-value")
	assert.NoError(t, err)
	assert.Equal(t, net.ParseIP("127.0.0.1"), ipValue)

	ipNetValue, err := flagSet.GetIPNet("ip-net-value")
	assert.NoError(t, err)
	assert.Equal(t, cfg.IPNetValue, ipNetValue)

	stringSliceValue, err := flagSet.GetStringSlice("string-slice-value")
	assert.NoError(t, err)
	assert.Equal(t, []string{"one", "two"}, stringSliceValue)

	intSliceValue, err := flagSet.GetIntSlice("int-slice-value")
	assert.NoError(t, err)
	assert.Equal(t, []int{10, 20}, intSliceValue)
}
