package main

import (
	"log"
	"os"
)

const minArgs = 3

func main() {
	if len(os.Args) < minArgs {
		log.Fatal("Usage: go-envdir /path/to/envdir command arg1 arg2")
	}

	dir, cmd := os.Args[1], os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	code := RunCmd(cmd, env)
	os.Exit(code)
}
