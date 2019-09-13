package commands

import (
	"github.com/nalej/connectivity-checker/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config = server.Config{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run connectivity-checker",
	Long:  `Run connectivity-checker`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Launching connectivity-checker!")
		server := server.NewService(config)
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
