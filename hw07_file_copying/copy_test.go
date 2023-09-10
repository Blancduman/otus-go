package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const bufferSize = 1024

var (
	testdataFolder = "./testdata/"
	readFileName   = "input.txt"
	writeFileName  = "out_offset%d_limit%d.txt"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		offset       int64
		limit        int64
		error        error
		readFileName string
	}{
		{
			offset:       0,
			limit:        0,
			error:        nil,
			readFileName: readFileName,
		},
		{
			offset:       0,
			limit:        1000,
			error:        nil,
			readFileName: readFileName,
		},
		{
			offset:       0,
			limit:        10000,
			error:        nil,
			readFileName: readFileName,
		},
		{
			offset:       100,
			limit:        1000,
			error:        nil,
			readFileName: readFileName,
		},
		{
			offset:       6000,
			limit:        1000,
			error:        nil,
			readFileName: readFileName,
		},
		{
			offset:       10000,
			limit:        1000,
			error:        ErrOffsetExceedsFileSize,
			readFileName: readFileName,
		},
		{
			offset:       0,
			limit:        0,
			error:        ErrUnsupportedFile,
			readFileName: "jkasjdklfjsakdlfjklsadjkf",
		},
	}

	err := os.Mkdir("./tmp", os.ModePerm)
	require.NoError(t, err)
	defer removeContents()

	for _, tt := range tests {
		tt := tt

		t.Run(fmt.Sprintf("limit %d offset %d", tt.limit, tt.offset), func(t *testing.T) {
			wfName := fmt.Sprintf(writeFileName, tt.offset, tt.limit)
			path := "./tmp/" + wfName

			err := Copy(testdataFolder+tt.readFileName, path, tt.offset, tt.limit)
			if tt.error != nil && errors.Is(err, tt.error) {
				return
			}

			require.NoError(t, err)

			testFile, err := os.Open(path)
			require.NoError(t, err)
			defer func() {
				err := testFile.Close()
				require.NoError(t, err)
			}()

			expectedFile, err := os.Open(testdataFolder + wfName)
			require.NoError(t, err)
			defer func() {
				err := expectedFile.Close()
				require.NoError(t, err)
			}()

			buffer1 := make([]byte, bufferSize)
			buffer2 := make([]byte, bufferSize)

			for {
				_, err1 := testFile.Read(buffer1)
				_, err2 := expectedFile.Read(buffer2)

				if !errors.Is(err1, io.EOF) {
					require.NoError(t, err1)
				}

				if !errors.Is(err2, io.EOF) {
					require.NoError(t, err2)
				}

				require.Equal(t, buffer2, buffer1)

				if err1 == io.EOF && err2 == io.EOF {
					break
				}
			}
		})
	}
}

func removeContents() {
	d, err := os.Open("./tmp")
	if err != nil {
		fmt.Print(err)
	}
	defer func() {
		err := os.Remove("./tmp/")
		if err != nil {
			fmt.Print(err)
		}
	}()
	defer func() {
		err := d.Close()
		if err != nil {
			fmt.Print(err)
		}
	}()
	names, err := d.Readdirnames(-1)
	if err != nil {
		fmt.Print(err)
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join("./tmp/", name))
		if err != nil {
			fmt.Print(err)
		}
	}
}
