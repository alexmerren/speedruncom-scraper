package repository

import (
	"encoding/csv"
	"os"
)

// AppendRepository creates a write-only append csv writer for a given
// filename. Use [AppendRepository] when you are processing data that requires
// multiple invocations of the same processor. i.e. Processing individual games
// by ID that are passed via command-line arguments.
type AppendRepository struct {
	writer *csv.Writer
}

// NewWriteRepository will create a [WriteRepository]. If the specified file
// does not exist, then create it. When the file exists, open it with mode
// [os.O_WRONLY] and [os.O_APPEND].
func NewAppendRepository(filename string) (*AppendRepository, func() error, error) {
	file, err := openOrCreate(filename, os.O_WRONLY|os.O_APPEND)
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

	return &AppendRepository{
		writer: writer,
	}, closeFunc, nil
}

func (w *AppendRepository) Write(record []string) error {
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
