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

type MinValidator struct {
	threshold int64
}

func (v MinValidator) Validate(val int64) error {
	if val < v.threshold {
		return ErrValueTooLow
	}
	return nil
}

type MaxValidator struct {
	threshold int64
}

func (v MaxValidator) Validate(val int64) error {
	if val > v.threshold {
		return ErrValueTooHigh
	}
	return nil
}

type IntInValidator struct {
	allowed []int64
}

func (v IntInValidator) Validate(val int64) error {
	if !slices.Contains(v.allowed, val) {
		return ErrValueNotInSet
	}
	return nil
}

func parseIntValidator(ruleType, ruleValue string) (Validator[int64], error) {
	var rule Validator[int64]
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

func parseMinRule(raw string) (Validator[int64], error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}

	return MinValidator{int64(v)}, nil
}

func parseMaxRule(raw string) (Validator[int64], error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}

	return MaxValidator{int64(v)}, nil
}

func parseIntInRule(raw string) (Validator[int64], error) {
	v, err := parseCommaSepInts(raw)
	if err != nil {
		return nil, err
	}

	return IntInValidator{v}, nil
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
