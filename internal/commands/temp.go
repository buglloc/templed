package commands

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/buglloc/templed/internal/sysinfo"
)

var tempArgs struct {
	Zone    int
	Monitor bool
	Period  time.Duration
}

var tempCmd = &cobra.Command{
	Use:           "temp",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Shows temp",
	RunE: func(_ *cobra.Command, _ []string) error {
		printTemp := func() error {
			temp, err := sysinfo.ZoneTemp(tempArgs.Zone)
			if err != nil {
				return err
			}

			fmt.Printf("%s\t%.2f\u2103\n", time.Now().Format("15:04:05"), temp)
			return nil
		}

		err := printTemp()
		if !tempArgs.Monitor || err != nil {
			return err
		}

		ticker := time.NewTicker(tempArgs.Period)
		defer ticker.Stop()

		for range ticker.C {
			if err := printTemp(); err != nil {
				log.Error().Err(err).Msg("unable to read temp")
				continue
			}
		}

		return nil
	},
}

func init() {
	flags := tempCmd.PersistentFlags()
	flags.IntVar(&tempArgs.Zone, "zone", 0, "thermal zone")
	flags.BoolVar(&tempArgs.Monitor, "monitor", false, "monitor temp")
	flags.DurationVar(&tempArgs.Period, "period", 1*time.Minute, "monitor period")
}
