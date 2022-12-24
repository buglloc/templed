package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:           "templed",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Templed is a daemon to adjust leds with temp",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(
		tempCmd,
		startCmd,
	)
}
