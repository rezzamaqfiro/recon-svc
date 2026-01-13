package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// Define flag for CLI
	totalRows := flag.Int("rows", 10000, "Total number of transactions to generate")
	flag.Parse()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	sysFile, _ := os.Create("system.csv")
	defer sysFile.Close()
	sysWriter := csv.NewWriter(sysFile)
	// Headers
	sysWriter.Write([]string{"TrxID", "Amount", "Type", "TransactionTime"})

	bankFiles := []*os.File{}
	bankWriters := []*csv.Writer{}
	bankNames := []string{"bank_a.csv", "bank_b.csv", "bank_c.csv"}

	for _, name := range bankNames {
		f, _ := os.Create(name)
		bankFiles = append(bankFiles, f)
		w := csv.NewWriter(f)
		// Headers
		w.Write([]string{"unique_identifier", "amount", "date"})
		bankWriters = append(bankWriters, w)
	}

	for i := 0; i < *totalRows; i++ {
		amount := int64(rand.Intn(1000000) + 1000)
		trxTime := startDate.Add(time.Duration(i) * time.Hour)

		sysWriter.Write([]string{
			fmt.Sprintf("TX%d", i),
			strconv.FormatInt(amount, 10),
			"CREDIT",
			trxTime.Format("2006-01-02 15:04:05"),
		})

		// variability in bank record
		dice := rand.Float64()
		if dice < 0.90 {
			// Randomly pick Bank A, B, or C
			bankIndex := rand.Intn(3)

			bankWriters[bankIndex].Write([]string{
				fmt.Sprintf("REF_%d", i),
				strconv.FormatInt(amount, 10),
				trxTime.Format("2006-01-02"),
			})
		}
	}

	sysWriter.Flush()
	for i, w := range bankWriters {
		w.Flush()
		bankFiles[i].Close()
	}

	fmt.Printf("\nSuccessfully generated  %d rows of test data system files (system.csv) and 3 bank files (bank_a.csv, bank_b.csv, bank_c.csv).\n", *totalRows)

}
