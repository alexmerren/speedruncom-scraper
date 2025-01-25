package repository

import (
	"encoding/csv"
	"os"
)

// ReadRepository creates a read-only csv reader for the given filename. Use
// when reading input files for processing data.
type ReadRepository struct {
	reader *csv.Reader
}

// NewReadRepository will create a [ReadRepository]. If the specified file does
// not exist then we throw an error. If the file does exist, open with
// [os.O_RDONLY] mode.
func NewReadRepository(filename string) (*ReadRepository, func() error, error) {
	file, err := openFile(filename, os.O_RDONLY)
	if err != nil {
		return nil, nil, err
	}

	csvReader := csv.NewReader(file)
	csvReader.Comment = '#'

	return &ReadRepository{
		reader: csvReader,
	}, file.Close, nil
}

func (r *ReadRepository) Read() ([]string, error) {
	return r.reader.Read()
}
