package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		if strings.ContainsRune(entry.Name(), '=') {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			env[entry.Name()] = EnvValue{"", true}
			continue
		}

		val, err := readFirstLineFrom(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		env[entry.Name()] = EnvValue{val, false}
	}

	return env, nil
}

func readFirstLineFrom(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	line = strings.TrimRight(line, " \t\n")
	line = strings.ReplaceAll(line, "\x00", "\n")

	return line, nil
}
