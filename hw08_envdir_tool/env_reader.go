package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var env Environment

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()

		f, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}

		defer func() {
			err := f.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		if env == nil {
			env = make(Environment)
		}

		scanner := bufio.NewScanner(f)
		if !scanner.Scan() {
			env[name] = EnvValue{NeedRemove: true}
			continue
		}

		value := strings.ReplaceAll(strings.TrimRight(scanner.Text(), " \t\n"), string(rune(0x00)), "\n")

		env[name] = EnvValue{
			Value:      value,
			NeedRemove: false,
		}
	}

	return env, nil
}
