package lib

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/mholt/archiver/v4"
)

type FileType = int

const (
	REGULAR FileType = iota
	DIRECTORY
	SYMLINK
)

type File struct {
	Path       string
	Mode       fs.FileMode
	LinkTarget string
	Type       FileType
}

func FileFromArchiver(f archiver.File, dest string) File {
	var fileType FileType

	if f.Mode().IsRegular() {
		fileType = REGULAR
	} else if f.Mode().IsDir() {
		fileType = DIRECTORY
	} else if f.LinkTarget != "" {
		fileType = SYMLINK
	}

	return File{
		Path:       dest,
		Mode:       f.Mode(),
		LinkTarget: f.LinkTarget,
		Type:       fileType,
	}
}

func FileFromOS(src, dest string) (File, error) {
	var fileType FileType

	fi, err := os.Lstat(src)
	if err != nil {
		return File{}, err
	}

	linkTarget := ""

	if fi.Mode().IsRegular() {
		fileType = REGULAR
	} else if fi.Mode().IsDir() {
		fileType = DIRECTORY
	} else if fi.Mode()&os.ModeType == os.ModeSymlink {
		link, err := os.Readlink(src)
		if err != nil {
			return File{}, err
		}
		linkTarget = link
		fileType = SYMLINK
	}

	return File{
		Path:       dest,
		Mode:       fi.Mode(),
		LinkTarget: linkTarget,
		Type:       fileType,
	}, nil
}

func CreateSymLink(f File) error {
	if f.LinkTarget == "" {
		return fmt.Errorf("symlink target is not specified")
	}
	return os.Symlink(f.LinkTarget, f.Path)
}

func CreateFile(f File, source func() io.ReadCloser) error {
	if !f.Mode.IsRegular() {
		return fmt.Errorf("given file is not a regular file")
	}

	destination, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source())
	return err
}

func CreateDirectory(f File) error {
	if !f.Mode.IsDir() {
		return fmt.Errorf("given file is not a directory")
	}

	if err := os.MkdirAll(f.Path, f.Mode.Perm()); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", f.Path, err.Error())
	}
	return nil
}

func SetPermissions(f File) error {
	return os.Chmod(f.Path, f.Mode.Perm())
}
