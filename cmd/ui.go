package cmd

import (
	"github.com/spf13/cobra"

	"machinelearning.one/sourcemap/compose/logger"
	"machinelearning.one/sourcemap/compose/server"
)

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Starts the web ui",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		lvl, _ := cmd.Flags().GetString("log-level")
		lg := logger.New(lvl)
		ctx = lg.WithContext(ctx)

		port, _ := cmd.Flags().GetUint("port")
		lg.Info().Uint("port", port).Msgf("Starting server on port %d", port)

		server.Run(ctx, port)
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().UintP("port", "p", 8080, "Port to listen on")
}