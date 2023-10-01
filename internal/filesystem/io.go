package filesystem

import (
	"errors"
	"os"
	"path/filepath"
)

func OpenInputFile(filename string) (*os.File, error) {
	return os.Open(filename)
}

func CreateOutputFile(filename string) (*os.File, error) {
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

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}
