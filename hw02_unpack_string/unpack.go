package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var builder strings.Builder
	var last string
	for _, cur := range input {
		if cur >= '0' && cur <= '9' {
			if last == "" {
				return "", ErrInvalidString
			}

			count, err := strconv.Atoi(string(cur))
			if err != nil {
				return "", err
			}
			builder.WriteString(strings.Repeat(last, count))
			last = ""
		} else {
			builder.WriteString(last)
			last = string(cur)
		}
	}

	if last != "" {
		builder.WriteString(last)
	}

	return builder.String(), nil
}
