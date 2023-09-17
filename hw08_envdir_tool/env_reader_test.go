package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) { //nolint: gocognit
	dir, err := os.MkdirTemp("./", "testdir")
	if err != nil {
		t.Fatalf("fail to create tmp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	}()

	t.Run("dir does not exist", func(t *testing.T) {
		env, err := ReadDir(filepath.Join(dir, "incorrentdir"))
		if err == nil {
			t.Fatal(errors.New("no error"))
		}

		if env != nil {
			t.Fatal("env is not nil")
		}
	})

	t.Run("file is not dir", func(t *testing.T) {
		fPath := filepath.Join(dir, "fatalbazooka")
		file, err := os.Create(fPath)
		if err != nil {
			t.Fatalf("fail to create file: %v", err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				t.Log(err)
			}

			if err := os.Remove(fPath); err != nil {
				t.Log(err)
			}
		}()

		_, err = file.Write([]byte("spaghetti"))
		if err != nil {
			t.Fatalf("fail to write: %v", err)
		}

		env, err := ReadDir(fPath)
		if err == nil {
			t.Fatal(errors.New("no error"))
		}

		if env != nil {
			t.Fatal(errors.New("env is not nil"))
		}
	})

	t.Run("empty dir", func(t *testing.T) {
		env, err := ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}

		if env != nil {
			t.Fatal(errors.New("env is not nil"))
		}
	})

	t.Run("files", func(t *testing.T) {
		err := os.WriteFile(filepath.Join(dir, "FOO"), []byte("123\n"), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		err = os.WriteFile(filepath.Join(dir, "BAR"), []byte("value"), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		env, err := ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}

		if len(env) != 2 {
			t.Fatal(errors.New("env len is not 2"))
		}

		if env["FOO"].Value != "123" {
			t.Fatal(errors.New("foo is not 123"))
		}

		if env["BAR"].Value != "value" {
			t.Fatal(errors.New("bar is not value"))
		}
	})
}
