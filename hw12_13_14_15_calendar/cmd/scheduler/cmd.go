package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var configFile string

func run(mCtx context.Context) error {
	ctx, cancel := signal.NotifyContext(mCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rootCmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage() //nolint:wrapcheck
		},
	}

	rootCmd.PersistentFlags().StringVarP(
		&configFile,
		"config",
		"c",
		"/etc/scheduler/scheduler_config.yaml",
		"Path to configuration file",
	)

	err := rootCmd.PersistentFlags().Parse(os.Args)
	if err != nil {
		return err
	}

	config := NewConfig(configFile)
	rootCmd.AddCommand(schedulerCmd(ctx, config))

	return errors.Wrap(rootCmd.ExecuteContext(ctx), "run scheduler")
}
