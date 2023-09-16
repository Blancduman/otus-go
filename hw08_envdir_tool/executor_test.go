package main

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()

	os.Stdout = w

	returnCode := RunCmd(
		[]string{"env"},
		Environment{
			"FOO": EnvValue{Value: "123"},
			"BAR": EnvValue{Value: "value"},
		},
	)

	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = oldStdout

	if returnCode != 0 {
		t.Fatal(errors.New("return code not 0"))
	}

	if !strings.Contains(string(output), "FOO=123") {
		t.Fatal(errors.New("no foo"))
	}

	if !strings.Contains(string(output), "BAR=value") {
		t.Fatal(errors.New("no bar"))
	}
}
