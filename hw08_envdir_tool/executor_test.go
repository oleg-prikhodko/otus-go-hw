package main

import "testing"

func TestRunCmd(t *testing.T) {
	t.Run("successful launch", func(t *testing.T) {
		code := RunCmd([]string{"pwd"}, nil)

		if code != 0 {
			t.Fatalf("return code is invalid: %v", code)
		}
	})

	t.Run("executable not found", func(t *testing.T) {
		code := RunCmd([]string{"234"}, nil)

		if code != 1 {
			t.Fatalf("return code is invalid: %v", code)
		}
	})

	t.Run("failed launch", func(t *testing.T) {
		code := RunCmd([]string{"pwd", "-x"}, nil)

		if code != 1 {
			t.Fatalf("return code is invalid: %v", code)
		}
	})
}
