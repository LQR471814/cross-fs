package archive

import (
	"github.com/spf13/cobra"
)

var Verbose *bool

var RootCmd = &cobra.Command{
	Use:   "archive",
	Short: "compress or extract an archive",
}
