package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	launch, args := cmd[0], cmd[1:]

	eCmd := exec.Command(launch, args...)
	eCmd.Stdout = os.Stdout
	eCmd.Stdin = os.Stdin
	eCmd.Stderr = os.Stderr

	if len(env) != 0 {
		for k, v := range env {
			if v.NeedRemove {
				err := os.Unsetenv(k)
				if err != nil {
					fmt.Println(err)
				}

				continue
			}

			if err := os.Setenv(k, v.Value); err != nil {
				fmt.Println(err)
			}
		}
	}

	if err := eCmd.Run(); err != nil {
		return 1
	}

	return 0
}
