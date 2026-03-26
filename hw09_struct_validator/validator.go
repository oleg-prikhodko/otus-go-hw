package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrIncorrectRuleString = errors.New("incorrect validation rule string")
	ErrNotStruct           = errors.New("value is not a struct")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%v: %v", v.Field, v.Err)
}

func (v ValidationError) Unwrap() error {
	return v.Err
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	messages := make([]string, len(v))
	for i, err := range v {
		messages[i] = err.Error()
	}

	return strings.Join(messages, "\n")
}

func (v ValidationErrors) Unwrap() []error {
	errs := make([]error, len(v))
	for i, e := range v {
		errs[i] = e
	}

	return errs
}

type Validator[T any] interface {
	Validate(val T) error
}

type RuleParser[T any] func(string, string) (Validator[T], error)

func parseValidators[T any](rawValidators string, parse RuleParser[T]) ([]Validator[T], error) {
	rawExpr := strings.Split(rawValidators, "|")

	validators := make([]Validator[T], 0)

	for _, e := range rawExpr {
		rule := strings.Split(e, ":")
		if len(rule) != 2 {
			return nil, ErrIncorrectRuleString
		}
		ruleType := rule[0]
		ruleValue := rule[1]

		r, err := parse(ruleType, ruleValue)
		if err != nil {
			return nil, err
		}
		validators = append(validators, r)
	}

	return validators, nil
}

func runValidators[T any](validators []Validator[T], name string, value T) ValidationErrors {
	var errs ValidationErrors

	for _, v := range validators {
		if err := v.Validate(value); err != nil {
			errs = append(errs, ValidationError{
				Field: name,
				Err:   err,
			})
		}
	}

	return errs
}

func forEach(v reflect.Value, fn func(int, any)) {
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		fn(i, elem.Interface())
	}
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	rt := rv.Type()

	var errs ValidationErrors

	for i := range rv.NumField() {
		fv := rv.Field(i)
		ft := rt.Field(i)
		if !ft.IsExported() {
			continue
		}

		rawValidators := ft.Tag.Get("validate")
		if rawValidators == "" {
			continue
		}

		switch fv.Kind() { //nolint:exhaustive
		case reflect.Int:
			validators, err := parseValidators(rawValidators, parseIntValidator)
			if err != nil {
				return err
			}
			errs = append(errs, runValidators(validators, ft.Name, fv.Int())...)
		case reflect.String:
			validators, err := parseValidators(rawValidators, parseStringValidator)
			if err != nil {
				return err
			}
			errs = append(errs, runValidators(validators, ft.Name, fv.String())...)
		case reflect.Slice:
			elemKind := ft.Type.Elem().Kind()

			if elemKind == reflect.Int {
				validators, err := parseValidators(rawValidators, parseIntValidator)
				if err != nil {
					return err
				}
				forEach(fv, func(i int, a any) {
					errs = append(errs, runValidators(validators, fmt.Sprintf("%v[%v]", ft.Name, i), a.(int64))...)
				})
			} else if elemKind == reflect.String {
				validators, err := parseValidators(rawValidators, parseStringValidator)
				if err != nil {
					return err
				}
				forEach(fv, func(i int, a any) {
					errs = append(errs, runValidators(validators, fmt.Sprintf("%v[%v]", ft.Name, i), a.(string))...)
				})
			}
		default:
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}
