package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/buglloc/templed/internal/config"
	"github.com/buglloc/templed/internal/watcher"
)

var startArgs struct {
	Configs []string
}

var startCmd = &cobra.Command{
	Use:           "start",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Starts daemon",
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := config.LoadConfig(startArgs.Configs...)
		if err != nil {
			return fmt.Errorf("unable to read config: %w", err)
		}

		instance, err := watcher.NewWatcher(cfg)
		if err != nil {
			return fmt.Errorf("unable to create watcher: %w", err)
		}

		errChan := make(chan error, 1)
		okChan := make(chan struct{})
		go func() {
			err := instance.Watch()
			if err != nil {
				errChan <- err
				return
			}

			close(okChan)
		}()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-stopChan:
			log.Info().Msg("shutting down gracefully by signal")

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			instance.Shutdown(ctx)
		case err := <-errChan:
			log.Error().Err(err).Msg("start failed")
			return err
		case <-okChan:
		}
		return nil
	},
}

func init() {
	flags := startCmd.PersistentFlags()
	flags.StringSliceVar(&startArgs.Configs, "cfg", nil, "cfg path")
}
