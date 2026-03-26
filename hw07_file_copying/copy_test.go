package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	input := "./testdata/input.txt"
	output := "./testdata/output.txt"

	t.Run("copy all bytes", func(t *testing.T) {
		err := Copy(input, output, 0, 0)
		if err != nil {
			t.Fatalf("err is not nil: %v", err)
		}
		defer cleanup(t, output)

		original := readFile(t, input)
		contents := readFile(t, output)

		if !bytes.Equal(original, contents) {
			t.Fatalf("incorrect output contents")
		}
	})

	t.Run("copy some bytes", func(t *testing.T) {
		err := Copy(input, output, 100, 200)
		if err != nil {
			t.Fatalf("err is not nil: %v", err)
		}
		defer cleanup(t, output)

		original := readFile(t, input)
		contents := readFile(t, output)

		if !bytes.Equal(original[100:300], contents) {
			t.Fatalf("incorrect output contents")
		}
	})

	t.Run("doesn't work on dirs", func(t *testing.T) {
		err := Copy("./testdata", output, 0, 0)
		if !errors.Is(err, ErrUnsupportedFile) {
			t.Fatalf("expecting unsupported file error, got: %v", err)
		}
	})

	t.Run("fails if offest too big", func(t *testing.T) {
		info, _ := os.Stat(input)

		err := Copy(input, output, info.Size()+1, 0)
		if !errors.Is(err, ErrOffsetExceedsFileSize) {
			t.Fatalf("expecting offset error, got: %v", err)
		}
	})
}

func cleanup(t *testing.T, filename string) {
	t.Helper()
	err := os.Remove(filename)
	if err != nil {
		t.Errorf("err is not nil: %v", err)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	return b
}
