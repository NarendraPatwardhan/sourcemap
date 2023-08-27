package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"machinelearning.one/sourcemap/compose/logger"
	"machinelearning.one/sourcemap/compose/server"
	"machinelearning.one/sourcemap/core"
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

		repo := core.Parse(
			ctx,
			"https://github.com/NarendraPatwardhan/sourcemap.git",
			core.ParseOpts{},
		)

		sm := server.API{
			HTTPMethod:   "GET",
			RelativePath: "/sourcemap",
			Handlers: []gin.HandlerFunc{
				func(c *gin.Context) {
					c.JSON(http.StatusOK, repo)
				},
			},
		}

		server.Run(ctx, port, sm)
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().UintP("port", "p", 8080, "Port to listen on")
}
