package hw09structvalidator

import (
	"errors"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedStringRule = errors.New("unsupported string validation rule type")
	ErrValueLength           = errors.New("value length incorrect")
	ErrValueNotCompatible    = errors.New("value is not compatible with regexp")
)

type LenValidator struct {
	length int
}

func (v LenValidator) Validate(val string) error {
	if len(val) != v.length {
		return ErrValueLength
	}
	return nil
}

type RegexpValidator struct {
	regexp *regexp.Regexp
}

func (v RegexpValidator) Validate(val string) error {
	if match := v.regexp.MatchString(val); !match {
		return ErrValueNotCompatible
	}
	return nil
}

type StringInValidator struct {
	allowed []string
}

func (v StringInValidator) Validate(val string) error {
	if !slices.Contains(v.allowed, val) {
		return ErrValueNotInSet
	}
	return nil
}

func parseStringValidator(ruleType, ruleValue string) (Validator[string], error) {
	var rule Validator[string]
	var err error

	switch ruleType {
	case "len":
		rule, err = parseLenRule(ruleValue)
	case "regexp":
		rule, err = parseRegexpRule(ruleValue)
	case "in":
		rule = parseStringInRule(ruleValue)
	default:
		return nil, ErrUnsupportedStringRule
	}

	if err != nil {
		return nil, err
	}

	return rule, nil
}

func parseLenRule(raw string) (Validator[string], error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}

	return LenValidator{v}, nil
}

func parseRegexpRule(raw string) (Validator[string], error) {
	v, err := regexp.Compile(raw)
	if err != nil {
		return nil, err
	}

	return RegexpValidator{v}, nil
}

func parseStringInRule(raw string) Validator[string] {
	v := strings.Split(raw, ",")

	return StringInValidator{v}
}
