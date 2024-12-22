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
		return file.Close()
	}

	return &WriteRepository{
		writer: writer,
	}, closeFunc, nil
}

func (w *WriteRepository) Write(record []string) error {
	return w.writer.Write(record)
}
