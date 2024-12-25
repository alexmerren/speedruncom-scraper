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

// openOrCreate will open a file if it already exists. If a file does not exist, then
//   - Create the file;
//   - Write the values from `FileComments` and `FileHeaders` (provided we have a
//     filename that matches a key in those maps);
//   - Open the file.
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

	if fileComment, ok := FileComments[filename]; ok {
		_, err = file.WriteString(fileComment + "\n")
		if err != nil {
			return nil, err
		}
	}

	if columnDefinitions, ok := FileColumnDefinitions[filename]; ok {
		_, err = file.WriteString(strings.Join(columnDefinitions, ",") + "\n")
		if err != nil {
			return nil, err
		}
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
