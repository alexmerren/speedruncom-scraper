package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func NewCsvReader(filename string) (*csv.Reader, func() error, error) {
	file, err := openFile(filename)
	if err != nil {
		return nil, nil, err
	}

	return csv.NewReader(file), file.Close, nil
}

func NewCsvWriter(filename string) (*csv.Writer, func() error, error) {
	file, err := openFile(filename)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(file)
	closeFunc := func() error {
		writer.Flush()
		return file.Close()
	}

	return writer, closeFunc, nil
}

// I have no idea where I stole this from. There's no way I wrote this.
func FormatCsvString(s string) string {
	s = fmt.Sprintf("%q", s)
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func openFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		directory, err := filepath.Abs(filepath.Dir(filename))
		if err != nil {
			return nil, err
		}

		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
}
