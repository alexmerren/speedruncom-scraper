package repository

import (
	"encoding/csv"
)

type WriteRepository struct {
	writer *csv.Writer
}

func NewWriteRepository(filename string) (*WriteRepository, func() error, error) {
	file, err := openOrCreate(filename)
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
