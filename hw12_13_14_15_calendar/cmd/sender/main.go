package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	err := run(context.Background())
	if err != nil {
		fmt.Println(errors.Wrap(err, "can not run command").Error())
		os.Exit(1)
	}

	os.Exit(0)
}
