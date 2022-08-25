package archive

import (
	"context"
	"crossfs/lib"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
	"github.com/spf13/cobra"
)

var unwrap *bool

var extractCmd = &cobra.Command{
	Use:   "extract [flags] file [destination]",
	Short: "extract an archive",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var input io.ReadCloser
		var err error

		input, err = os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}

		format, _, err := archiver.Identify(args[0], input)
		if err != nil {
			log.Fatal(err)
		}

		if decompress, ok := format.(archiver.Decompressor); ok {
			rc, err := decompress.OpenReader(input)
			if err != nil {
				log.Fatal(err)
			}
			defer rc.Close()
			input = rc
		}

		if ex, ok := format.(archiver.Extractor); ok {
			unwrappable := true
			rootDirectories := 0
			rootPath := ""

			err := ex.Extract(
				context.Background(), input, nil,
				func(ctx context.Context, f archiver.File) error {
					name := filepath.Join(args[1], f.NameInArchive)
					file := lib.FileFromArchiver(f, name)

					if *unwrap {
						//* strings.Split() use instead of filepath.SplitList() because
						//* archiver name will always be separated with slashes
						path := strings.Split(f.NameInArchive, "/")
						filtered := []string{}
						for _, s := range path {
							if strings.Trim(s, " ") != "" {
								filtered = append(filtered, s)
							}
						}
						if len(filtered) == 1 {
							if f.IsDir() {
								rootDirectories++
								rootPath = name
							} else {
								unwrappable = false
							}
							if rootDirectories > 1 {
								unwrappable = false
							}
						}
					}

					info := func(file string) {
						if *Verbose {
							log.Printf("[extract] %s %s", file, name)
						}
					}

					switch file.Type {
					case lib.REGULAR:
						info("regular")
						return lib.CreateFile(file, func() io.ReadCloser {
							source, err := f.Open()
							if err != nil {
								log.Fatal(err)
							}
							return source
						})
					case lib.DIRECTORY:
						info("directory")
						return lib.CreateDirectory(file)
					case lib.SYMLINK:
						info("symlink")
						return lib.CreateSymLink(file)
					}

					return nil
				},
			)
			if err != nil {
				log.Fatal(err)
			}

			if *unwrap && unwrappable {
				if *Verbose {
					log.Printf("[extract] unwrap %s -> %s \n", rootPath, args[1])
				}

				temporary := filepath.Join(
					filepath.Dir(args[1]),
					fmt.Sprintf("../%d", rand.Int63()),
				)
				err := os.Rename(rootPath, temporary)
				if err != nil {
					log.Fatal(err)
				}

				err = os.Remove(args[1])
				if err != nil {
					log.Fatal(err)
				}

				err = os.Rename(temporary, args[1])
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	},
}

func init() {
	unwrap = extractCmd.Flags().BoolP(
		"unwrap", "u", false,
		"replace the outermost directory of the archive with that of the extract destination",
	)
	RootCmd.AddCommand(extractCmd)
}
