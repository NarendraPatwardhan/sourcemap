package cmd

import (
	"github.com/spf13/cobra"

	cctx "machinelearning.one/sourcemap/compose/context"
	"machinelearning.one/sourcemap/compose/logger"
	"machinelearning.one/sourcemap/core"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sourcemap",
	Short: "A git repository visualizer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr := args[0]
		ctx := cmd.Context()
		lvl, _ := cmd.Flags().GetString("log-level")
		lg := logger.New(lvl)
		ctx = lg.WithContext(ctx)

		limit, _ := cmd.Flags().GetInt("limit")
		excludeGlobs, _ := cmd.Flags().GetStringSlice("exclude-globs")
		excludePaths, _ := cmd.Flags().GetStringSlice("exclude-paths")

		lg.Trace().Msgf("using parse limit %d", limit)
		lg.Trace().Msgf("using exclude globs %v", excludeGlobs)
		lg.Trace().Msgf("using exclude paths %v", excludePaths)

		lg.Debug().Msgf("parsing %s", addr)
		core.Parse(ctx, addr, core.ParseOpts{
			Limit:        limit,
			ExcludeGlobs: excludeGlobs,
			ExcludePaths: excludePaths,
		})
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	ctx := cctx.Context()
	rootCmd.SetContext(ctx)

	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("log-level", "l", logger.DefaultLevel, "log level")
	rootCmd.Flags().IntP("limit", "n", 0, "the number of commits to parse, (default \"all\")")
	rootCmd.Flags().
		StringSliceP("exclude-globs", "g", []string{}, "exclude files matching the given glob patterns")
	rootCmd.Flags().
		StringSliceP("exclude-paths", "p", []string{}, "exclude files matching the given paths")
}
