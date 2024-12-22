package repository

import "encoding/csv"

type ReadRepository struct {
	reader *csv.Reader
}

func NewReadRepository(filename string) (*ReadRepository, func() error, error) {
	file, err := openFile(filename)
	if err != nil {
		return nil, nil, err
	}

	return &ReadRepository{
		reader: csv.NewReader(file),
	}, file.Close, nil
}

func (r *ReadRepository) Read() ([]string, error) {
	return r.reader.Read()
}
