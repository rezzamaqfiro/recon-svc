package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ParseSystemCSV(filePath string) ([]Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var transactions []Transaction
	// TODO: Read and parse CSV rows into Transaction structs
	return transactions, nil
}
