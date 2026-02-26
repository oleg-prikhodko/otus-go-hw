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

type StringRule int

const (
	Len StringRule = iota
	Regexp
	StringIn
)

type StringValidationRule struct {
	ruleType StringRule
	val      any
}

func (r StringValidationRule) Validate(val string) error {
	switch r.ruleType {
	case Len:
		if len(val) != r.val.(int) {
			return ErrValueLength
		}
	case Regexp:
		if match := r.val.(*regexp.Regexp).MatchString(val); !match {
			return ErrValueNotCompatible
		}
	case StringIn:
		if !slices.Contains(r.val.([]string), val) {
			return ErrValueNotInSet
		}
	}

	return nil
}

func parseStringValidator(ruleType, ruleValue string) (Validator[string], error) {
	var rule StringValidationRule
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

func parseLenRule(raw string) (StringValidationRule, error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return StringValidationRule{}, err
	}

	return StringValidationRule{Len, v}, nil
}

func parseRegexpRule(raw string) (StringValidationRule, error) {
	v, err := regexp.Compile(raw)
	if err != nil {
		return StringValidationRule{}, err
	}

	return StringValidationRule{Regexp, v}, nil
}

func parseStringInRule(raw string) StringValidationRule {
	v := strings.Split(raw, ",")

	return StringValidationRule{StringIn, v}
}
