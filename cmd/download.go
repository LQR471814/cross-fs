package cmd

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gobwas/httphead"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download [flags] url [destination]",
	Short: "downloads a file from an http server",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{}
		response, err := client.Get(args[0])
		if err != nil {
			log.Fatal(err)
		}

		destination := "downloaded-file"
		if len(args) > 1 {
			destination = args[1]
		} else {
			options, ok := httphead.ParseOptions(
				[]byte(response.Header.Get("Content-Disposition")), nil,
			)
			if ok {
				for _, o := range options {
					if string(o.Name) == "attachment" {
						value, ok := o.Parameters.Get("filename")
						if !ok {
							value, ok = o.Parameters.Get("filename*")
							if !ok {
								continue
							}
						}
						destination = string(value)
						break
					}
				}
			}
		}

		f, err := os.Create(destination)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if *verbose {
			log.Printf("[download] %s -> %s\n", args[0], destination)
		}

		_, err = io.Copy(f, response.Body)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
