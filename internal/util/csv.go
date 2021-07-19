package util

import (
	"encoding/csv"
	"os"
)

func ReadCsv(file string) ([][]string, error) {
	f, err := os.Open(file)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
