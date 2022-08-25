package lib

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CopyOptions struct {
	Source      string
	Destination string

	Verbose bool
}

func Copy(o CopyOptions) error {
	info := func(mode string) {
		if o.Verbose {
			prefix := "copy"
			log.Printf(
				"[%s] %-4s %s -> %s\n",
				prefix, strings.ToUpper(mode),
				o.Source, o.Destination,
			)
		}
	}

	f, err := FileFromOS(o.Source, o.Destination)
	if err != nil {
		return err
	}

	switch f.Type {
	case REGULAR:
		info("file")
		err := CreateFile(f, func() io.ReadCloser {
			source, err := os.Open(o.Source)
			if err != nil {
				log.Fatal(err)
			}
			return source
		})
		if err != nil {
			return err
		}
	case DIRECTORY:
		info("dir")
		err := CreateDirectory(f)
		if err != nil {
			return err
		}
		err = copyDirectory(o)
		if err != nil {
			return err
		}
	case SYMLINK:
		info("link")
		err := CreateSymLink(f)
		if err != nil {
			return err
		}
	}

	if o.Verbose {
		log.Printf("[chmod] %s %v", f.Path, f.Mode.Perm())
	}
	return SetPermissions(f)
}

func copyDirectory(o CopyOptions) error {
	entries, err := os.ReadDir(o.Source)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		entryOptions := o
		entryOptions.Source = filepath.Join(o.Source, entry.Name())
		entryOptions.Destination = filepath.Join(o.Destination, entry.Name())

		err := Copy(entryOptions)
		if err != nil {
			return err
		}
	}
	return nil
}
