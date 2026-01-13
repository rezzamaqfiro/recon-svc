package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
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
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		amtFloat, _ := strconv.ParseFloat(record[1], 64)
		amtIDR := int64(amtFloat)

		// parse time
		t, _ := time.Parse("2006-01-02 15:04:05", record[3])

		transactions = append(transactions, Transaction{
			TrxID:  record[0],
			Amount: amtIDR,
			Type:   record[2],
			Time:   t,
		})
	}
	return transactions, nil
}
