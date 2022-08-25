package cmd

import (
	"crossfs/cmd/archive"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var verbose *bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fs",
	Short: "cross platform filesystem operations",
}

func GenerateDocs(dir string) error {
	return doc.GenMarkdownTree(rootCmd, dir)
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	verbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")

	archive.Verbose = verbose
	rootCmd.AddCommand(archive.RootCmd)
}
