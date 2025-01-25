package repository

import (
	"encoding/csv"
	"os"
)

// WriteRepository creates a write-only csv writer for a given filename. Use
// [WriteRepository] when you are processing data that is only invocated once,
// as the file will flush every time. i.e. Processing all games from an input
// file. If you are processing data and the file header is missing, it is
// likely due to invocating a write repository on an already existing file.
type WriteRepository struct {
	writer *csv.Writer
}

// NewWriteRepository will create a [WriteRepository]. If the specified file
// does not exist, then create it. When the file exists, open it with mode
// [os.O_WRONLY].
func NewWriteRepository(filename string) (*WriteRepository, func() error, error) {
	file, err := openOrCreate(filename, os.O_WRONLY)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(file)

	closeFunc := func() error {
		writer.Flush()
		err := writer.Error()
		if err != nil {
			return err
		}

		return file.Close()
	}

	return &WriteRepository{
		writer: writer,
	}, closeFunc, nil
}

func (w *WriteRepository) Write(record []string) error {
	err := w.writer.Write(record)
	if err != nil {
		return err
	}

	err = w.writer.Error()
	if err != nil {
		return err
	}

	return nil
}
