package cmd

import (
	"io"
	"os"
)

// fileSystem is the interface that simulates the os file system
type fileSystem interface {
	Create(name string) (file, error)
}

// file is the interface which simulates os.File
type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	Stat() (os.FileInfo, error)
}
