package hw09structvalidator

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedIntRule = errors.New("unsupported int validation rule type")
	ErrValueTooLow        = errors.New("value too low")
	ErrValueTooHigh       = errors.New("value too high")
	ErrValueNotInSet      = errors.New("value is not in the set")
)

type IntRule int

const (
	Min IntRule = iota
	Max
	IntIn
)

type IntValidationRule struct {
	ruleType IntRule
	val      any
}

func (r IntValidationRule) Validate(val int64) error {
	switch r.ruleType {
	case Min:
		if val < r.val.(int64) {
			return ErrValueTooLow
		}
	case Max:
		if val > r.val.(int64) {
			return ErrValueTooHigh
		}
	case IntIn:
		if !slices.Contains(r.val.([]int64), val) {
			return ErrValueNotInSet
		}
	}

	return nil
}

func parseIntValidator(ruleType, ruleValue string) (Validator[int64], error) {
	var rule IntValidationRule
	var err error

	switch ruleType {
	case "min":
		rule, err = parseMinRule(ruleValue)
	case "max":
		rule, err = parseMaxRule(ruleValue)
	case "in":
		rule, err = parseIntInRule(ruleValue)
	default:
		return nil, ErrUnsupportedIntRule
	}

	if err != nil {
		return nil, err
	}

	return rule, nil
}

func parseMinRule(raw string) (IntValidationRule, error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return IntValidationRule{}, err
	}

	return IntValidationRule{Min, int64(v)}, nil
}

func parseMaxRule(raw string) (IntValidationRule, error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return IntValidationRule{}, err
	}

	return IntValidationRule{Max, int64(v)}, nil
}

func parseIntInRule(raw string) (IntValidationRule, error) {
	v, err := parseCommaSepInts(raw)
	if err != nil {
		return IntValidationRule{}, err
	}

	return IntValidationRule{IntIn, v}, nil
}

func parseCommaSepInts(raw string) ([]int64, error) {
	vals := make([]int64, 0)
	for _, v := range strings.Split(raw, ",") {
		iv, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, int64(iv))
	}

	return vals, nil
}
