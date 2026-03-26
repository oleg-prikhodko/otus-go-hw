package main

import (
	"testing"
)

func TestReadDir(t *testing.T) {
	env, err := ReadDir("./testdata/env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := Environment{
		"BAR":   EnvValue{"bar", false},
		"EMPTY": EnvValue{"", false},
		"FOO":   EnvValue{"   foo\nwith new line", false},
		"HELLO": EnvValue{"\"hello\"", false},
		"UNSET": {"", true},
	}

	if len(expected) != len(env) {
		t.Fatalf("expected map len %v, got %v", len(expected), len(env))
	}
	for k, v := range env {
		if eVal, ok := expected[k]; !ok || eVal != v {
			t.Fatalf("expected %v val %v, but got %v", k, eVal, v)
		}
	}
}
