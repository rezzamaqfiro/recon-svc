package main

import (
	"fmt"
)

func main() {
	sysData, err := ParseSystemCSV("internal_trxs.csv")
	if err != nil {
		fmt.Printf("Error loading system data: %v\n", err)
		return
	}

	bankData, _ := ParseBankCSV("bank_records.csv")
	summary := Reconcile(sysData, bankData)

	fmt.Printf("--- Reconciliation Summary ---\n")
	fmt.Printf("Total Processed: %d\n", summary.TotalProcessed)
	fmt.Printf("Matched Count: %d\n", summary.MatchedCount)
	fmt.Printf("Total Discrepancy: %d\n", summary.TotalDiscrepancy)

	if summary.TotalDiscrepancy > 0 {
		fmt.Printf("Total Amount Gap: Rp %.2f", float64(summary.TotalDiscrepancy)/100)
	}
}
