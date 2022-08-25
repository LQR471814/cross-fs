package cmd

import (
	"crossfs/cmd/archive"

	"github.com/spf13/cobra"
)

var verbose *bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fs",
	Short: "cross platform filesystem operations",
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	verbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")

	archive.Verbose = verbose
	rootCmd.AddCommand(archive.RootCmd)
}
