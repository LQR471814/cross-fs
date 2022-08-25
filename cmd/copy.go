package cmd

import (
	"crossfs/lib"
	"log"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy [flags] src dest",
	Short: "copy a filesystem entry",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := lib.Copy(lib.CopyOptions{
			Source:      args[0],
			Destination: args[1],
			Verbose:     *verbose,
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
