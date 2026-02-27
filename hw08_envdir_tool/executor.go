package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	osEnv := updateOsEnv(env)
	for k, v := range osEnv {
		c.Env = append(c.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if err := c.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}

func updateOsEnv(env Environment) map[string]string {
	osEnv := envToMap()
	for k, v := range env {
		if v.NeedRemove {
			delete(osEnv, k)
		} else {
			osEnv[k] = v.Value
		}
	}

	return osEnv
}

func envToMap() map[string]string {
	envMap := make(map[string]string)
	for _, v := range os.Environ() {
		parts := strings.SplitN(v, "=", 2)
		envMap[parts[0]] = parts[1]
	}

	return envMap
}
