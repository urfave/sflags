package govalidator

import (
	"flag"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/octago/sflags"
	"github.com/octago/sflags/gen/gflag"
	"github.com/stretchr/testify/assert"
)

func ExampleNew() {
	type config struct {
		Host string `valid:"host"`
		Port int    `valid:"port"`
	}
	cfg := &config{
		Host: "127.0.0.1",
		Port: 6000,
	}
	// Use gflags.ParseToDef if you want default `flag.CommandLine`
	fs, err := gflag.Parse(cfg, sflags.Validator(New()))
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	fs.Init("text", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	// if we pass a wrong domain to the host flag, we'll get a error.
	if err = fs.Parse([]string{"-host", "wrong domain"}); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	// if we pass a wrong port to the port flag, we'll get a error.
	if err = fs.Parse([]string{"-port", "800000"}); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	// Output:
	// err: invalid value "wrong domain" for flag -host: `wrong domain` does not validate as host
	// err: invalid value "800000" for flag -port: `800000` does not validate as port
}

func Test_isValidTag(t *testing.T) {
	tests := []struct {
		arg  string
		want bool
	}{
		{"simple", true},
		{"", false},
		{"!#$%&()*+-./:<=>?@[]^_{|}~ ", true},
		{"абв", true},
		{"`", false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, isValidTag(tt.arg), "for %v", tt.arg)
	}
}

func Test_parseTagIntoMap(t *testing.T) {
	tests := []struct {
		tag  string
		want tagOptionsMap
	}{
		{
			tag: "required~Some error message,length(2|3)",
			want: tagOptionsMap{
				"required":    "Some error message",
				"length(2|3)": "",
			},
		},
		{
			tag: "required~Some error message~other",
			want: tagOptionsMap{
				"required": "",
			},
		},
		{
			tag: "bad`tag,good_tag",
			want: tagOptionsMap{
				"good_tag": "",
			},
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, parseTagIntoMap(tt.tag), "for %v", tt.tag)
	}
}

func Test_validateFunc(t *testing.T) {
	tests := []struct {
		val     string
		options tagOptionsMap

		expErr string
	}{
		{
			val:     "not a host",
			options: tagOptionsMap{"host": ""},
			expErr:  "`not a host` does not validate as host",
		},
		{
			val:     "localhost",
			options: tagOptionsMap{"host": ""},
			expErr:  "",
		},
		{
			val:     "localhost",
			options: tagOptionsMap{"!host": ""},
			expErr:  "`localhost` does validate as host",
		},
		{
			val:     "not a host",
			options: tagOptionsMap{"host": "wrong host value"},
			expErr:  "wrong host value",
		},
		{
			val:     "localhost",
			options: tagOptionsMap{"!host": "shouldn't be a host"},
			expErr:  "shouldn't be a host",
		},
		{
			val:     "localhost",
			options: tagOptionsMap{"length(2|10)": ""},
			expErr:  "",
		},
		{
			val:     "localhostlong",
			options: tagOptionsMap{"length(2|10)": ""},
			expErr:  "`localhostlong` does not validate as length(2|10)",
		},
		{
			val:     "localhostlong",
			options: tagOptionsMap{"length(2|10)": "too long!"},
			expErr:  "too long!",
		},
		{
			val:     "localhost",
			options: tagOptionsMap{"!length(2|10)": ""},
			expErr:  "`localhost` does validate as length(2|10)",
		},
		{
			val:     "localhost",
			options: tagOptionsMap{"!length(2|10)": "should be longer"},
			expErr:  "should be longer",
		},
	}
	for _, tt := range tests {
		err := validateFunc(tt.val, tt.options)
		if tt.expErr != "" {
			if assert.Error(t, err) {
				assert.EqualError(t, err, tt.expErr)
			}
		} else {
			assert.NoError(t, err)
		}
	}
}
