package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FormatCsvString(s string) string {
	s = fmt.Sprintf("%q", s)
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// openOrCreate will open a file if it already exists. If a file does not exist,
// then we create the file and write the specified values from `FileComments` and
// `FileHeaders` to the file, provided we have a filename that matches a key in
// those maps.
func openOrCreate(filename string) (*os.File, error) {
	err := createFileIfNotExists(filename)
	if err != nil && errors.Is(err, os.ErrExist) {
		return openFile(filename)
	}

	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}

	file, err := openFile(filename)
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(FileComments[filename])
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(strings.Join(FileColumnDefinitions[filename], ","))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func createFileIfNotExists(filename string) error {
	_, err := os.Stat(filename)
	if err == nil {
		return os.ErrExist
	}

	if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	directory, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		return err
	}

	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(filename)

	return err
}

func openFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR, 0600)
}
