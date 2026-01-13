package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// Define flag for CLI
	sysPath := flag.String("system", "", "Path to system transaction CSV")
	bankPath := flag.String("bank", "", "Path to bank statement CSV")
	windowDays := flag.Int("window", 1, "Time window in days (optional, default 1)")
	flag.Parse()

	if *sysPath == "" || *bankPath == "" {
		log.Fatal("Please provide both -system and -bank file paths ")
	}

	sysTrxs, err := ParseSystemCSV(*sysPath)
	if err != nil {
		log.Fatalf("System parse error: %v", err)
	}

	bankRecs, err := ParseBankCSV(*bankPath)
	if err != nil {
		log.Fatalf("Bank parse error: %v", err)
	}

	summary := Reconcile(sysTrxs, bankRecs, *windowDays)
	fmt.Printf("\n--- AMARTHA RECONCILIATION REPORT ---\n")
	fmt.Printf("Matched : %d\n", summary.MatchedCount)
	fmt.Printf("Unmatched : %d\n", summary.UnmatchedCount)
	fmt.Printf("Discrepancy Total: Rp.%d,-\n", summary.TotalDiscrepancy)
}
