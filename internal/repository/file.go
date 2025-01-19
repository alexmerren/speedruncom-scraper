package repository

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// openOrCreate will open a file if it already exists. If a file does not exist, then
//   - Create the file;
//   - Write the value from `FileHeaders` (provided we have a
//     filename that matches a key in the map);
//   - Open the file.
func openOrCreate(filename string, mode int) (*os.File, error) {
	fileExists, err := doesFileExist(filename)
	if err != nil {
		return nil, err
	}

	if fileExists {
		slog.Info("Opening existing file", "filename", filename)
		return openFile(filename, mode)
	}

	slog.Info("Creating missing file", "filename", filename)

	file, err := createFile(filename)
	if err != nil {
		return nil, err
	}

	if columnDefinitions, ok := FileColumnDefinitions[filename]; ok {
		_, err = file.WriteString(strings.Join(columnDefinitions, ",") + "\n")
		if err != nil {
			return nil, err
		}
	}

	err = file.Sync()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func doesFileExist(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	return true, nil
}

func createFile(filename string) (*os.File, error) {
	directory, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return os.Create(filename)
}

func openFile(filename string, mode int) (*os.File, error) {
	return os.OpenFile(filename, mode, 0600)
}
