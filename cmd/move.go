package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move [flags] src dest",
	Short: "move a filesystem entry",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Rename(args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
