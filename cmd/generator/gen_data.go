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
	bankFile, _ := os.Create("bank.csv")
	defer sysFile.Close()
	defer bankFile.Close()

	sysWriter := csv.NewWriter(sysFile)
	bankWriter := csv.NewWriter(bankFile)
	defer sysWriter.Flush()
	defer bankWriter.Flush()

	// Headers
	sysWriter.Write([]string{"TrxID", "Amount", "Type", "TransactionTime"})
	bankWriter.Write([]string{"UID", "Amount", "Date"})

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
		if dice < 0.85 {
			settleDate := trxTime.AddDate(0, 0, rand.Intn(2))
			bankWriter.Write([]string{
				fmt.Sprintf("BANK_%d", i),
				strconv.FormatInt(amount, 10),
				settleDate.Format("2006-01-02"),
			})
		} else if dice < 0.95 {
			diff := int64(rand.Intn(500) + 1)
			bankWriter.Write([]string{
				fmt.Sprintf("BANK_ERR_%d", i),
				strconv.FormatInt(amount+diff, 10),
				trxTime.Format("2006-01-02"),
			})
		} else {
			// skip to simulate missing record
		}
	}

	sysWriter.Flush()
	bankWriter.Flush()
	fmt.Printf("\nGenerated %d rows of test data successfully.\n", *totalRows)

}
