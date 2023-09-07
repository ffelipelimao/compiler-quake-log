package reader

import (
	"bufio"
	"fmt"
	"os"
)

type Reader struct {
	Path string
}

func New(path string) *Reader {
	return &Reader{
		Path: path,
	}
}

func (c *Reader) LoadRows() ([]string, error) {
	var rows []string

	file, err := os.Open(c.Path)
	if err != nil {
		return nil, fmt.Errorf("error to open file %s", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		rows = append(rows, row)
	}

	return rows, nil
}
