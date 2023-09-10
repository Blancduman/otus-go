package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	readFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	defer func() {
		err := readFile.Close()
		if err != nil {
			fmt.Printf("can not close read file file: %s", err)
		}
	}()

	writeFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("can not close created file: %w", err)
	}

	defer func() {
		err := writeFile.Close()
		if err != nil {
			fmt.Printf("can not close read file file: %s", err)
		}
	}()

	stats, err := readFile.Stat()
	if err != nil {
		return fmt.Errorf("can not read file stats: %w", err)
	}

	readFileSize := stats.Size()

	if offset > readFileSize {
		return ErrOffsetExceedsFileSize
	}

	var readFileLimit int64

	if limit == 0 || offset+limit > readFileSize {
		readFileLimit = readFileSize - offset
	} else {
		readFileLimit = limit
	}

	_, err = readFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("can not seek read file: %w", err)
	}

	pBar := progressbar.Default(readFileLimit, "copying")
	defer func() {
		err := pBar.Close()
		if err != nil {
			fmt.Printf("can not close progressbar: %s", err)
		}
	}()

	writeFileProgress := io.MultiWriter(pBar, writeFile)

	_, err = io.CopyN(writeFileProgress, readFile, readFileLimit)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("can not copy file: %w", err)
	}

	err = pBar.Finish()
	if err != nil {
		return fmt.Errorf("can not finish progressbar: %w", err)
	}

	return nil
}
