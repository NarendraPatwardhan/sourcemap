package cmd

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

		godotenv.Load()

		port, _ := cmd.Flags().GetUint("port")
		decoupled, _ := cmd.Flags().GetBool("decoupled")

		mode := "integrated"
		if decoupled {
			mode = "decoupled"
		}

		if decoupled && port != 0 {
			lg.Fatal().
				Msg("Cannot specify port manually when running in decoupled mode, use SOURCEMAP_API_PORT instead")
		}
		if port == 0 {
			portString := os.Getenv("SOURCEMAP_API_PORT")
			parsed, err := strconv.Atoi(portString)
			if err != nil {
				lg.Fatal().Err(err).Msg("Could not parse SOURCEMAP_API_PORT")
			} else {
				port = uint(parsed)
			}
		}

		repo := core.Parse(
			ctx,
			"https://github.com/NarendraPatwardhan/sourcemap.git",
			core.ParseOpts{},
		)

		sm := server.Func{
			HTTPMethod:   "GET",
			RelativePath: "/sourcemap",
			Handlers: []gin.HandlerFunc{
				func(c *gin.Context) {
					c.JSON(http.StatusOK, repo)
				},
			},
		}

		lg.Info().Msgf("Starting server on port %d in %s mode", port, mode)
		server.Run(ctx, port, decoupled, sm)
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().UintP("port", "p", 0, "Port to listen on")
	uiCmd.Flags().BoolP("decoupled", "d", false, "Whether to run API server only")
}
