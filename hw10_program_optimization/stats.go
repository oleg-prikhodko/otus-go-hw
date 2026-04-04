package hw10programoptimization

import (
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type PartialUser struct {
	Email string
}

type users [100_000]User

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	suffix := "." + domain
	result := make(DomainStat)
	dec := json.NewDecoder(r)
	var user PartialUser

	for dec.More() {
		if err := dec.Decode(&user); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, suffix) {
			idx := strings.IndexRune(user.Email, '@')
			result[strings.ToLower(user.Email[idx+1:])]++
		}
	}

	return result, nil
}
