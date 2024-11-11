// Package govalidator adds support for govalidator library.
// Part of this package was taken from govalidator private api
// and covered by MIT license.
//
// The MIT License (MIT)
//
// # Copyright (c) 2014 Alex Saskevich
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package govalidator

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/asaskevich/govalidator"
)

const (
	validTag = "valid" // tag isn't optional in govalidator
)

type tagOptionsMap map[string]string

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
		// Backslash and quote chars are reserved, but
		// otherwise any punctuation chars are allowed
		// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}

// parseTagIntoMap parses a struct tag `valid:"required~Some error message,length(2|3)"` into map[string]string{"required": "Some error message", "length(2|3)": ""}
func parseTagIntoMap(tag string) tagOptionsMap {
	optionsMap := make(tagOptionsMap)
	options := strings.SplitN(tag, ",", -1)
	for _, option := range options {
		validationOptions := strings.Split(option, "~")
		if !isValidTag(validationOptions[0]) {
			continue
		}
		if len(validationOptions) == 2 {
			optionsMap[validationOptions[0]] = validationOptions[1]
		} else {
			optionsMap[validationOptions[0]] = ""
		}
	}
	return optionsMap
}

func validateFunc(val string, options tagOptionsMap) error {

	// for each tag option check the map of validator functions
	for validator, customErrorMessage := range options {
		var negate bool
		customMsgExists := len(customErrorMessage) > 0
		// Check wether the tag looks like '!something' or 'something'
		if validator[0] == '!' {
			validator = string(validator[1:])
			negate = true
		}
		// Check for param validators
		for key, value := range govalidator.ParamTagRegexMap {
			ps := value.FindStringSubmatch(validator)
			if len(ps) > 0 {
				if validatefunc, ok := govalidator.ParamTagMap[key]; ok {

					if result := validatefunc(val, ps[1:]...); (!result && !negate) || (result && negate) {
						var err error
						if !negate {
							if customMsgExists {
								err = fmt.Errorf("%s", customErrorMessage)
							} else {
								err = fmt.Errorf("`%s` does not validate as %s", val, validator)
							}

						} else {
							if customMsgExists {
								err = fmt.Errorf("%s", customErrorMessage)
							} else {
								err = fmt.Errorf("`%s` does validate as %s", val, validator)
							}
						}
						return err
					}

				}
			}
		}

		if validatefunc, ok := govalidator.TagMap[validator]; ok {
			if result := validatefunc(val); !result && !negate || result && negate {
				var err error

				if !negate {
					if customMsgExists {
						err = fmt.Errorf("%s", customErrorMessage)
					} else {
						err = fmt.Errorf("`%s` does not validate as %s", val, validator)
					}
				} else {
					if customMsgExists {
						err = fmt.Errorf("%s", customErrorMessage)
					} else {
						err = fmt.Errorf("`%s` does validate as %s", val, validator)
					}
				}
				return err
			}

		}
	}
	return nil
}

// New returns ValidateFunc for govalidator library.
// Supports default String validators in TagMap and ParamTagMap.
// Doesn't support custom type validators and required filters.
// Please check all available functions at https://github.com/asaskevich/govalidator.
func New() func(val string, field reflect.StructField, obj interface{}) error {
	return func(val string, field reflect.StructField, obj interface{}) error {
		options := parseTagIntoMap(field.Tag.Get(validTag))
		return validateFunc(val, options)
	}
}
