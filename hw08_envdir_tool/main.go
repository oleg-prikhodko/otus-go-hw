package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("invalid number of arguments")
		os.Exit(1)
	}

	dir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println("error while reading dir")
		os.Exit(1)
	}

	code := RunCmd(cmd, env)
	os.Exit(code)
}
