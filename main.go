package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	// Define flag for CLI
	sysPath := flag.String("system", "", "Path to system transaction CSV")
	bankPaths := flag.String("bank", "", "Path to bank statement CSV")
	windowDays := flag.Int("window", 1, "Time window in days (optional, default 1)")
	startStr := flag.String("start", "", "Start date (YYYY-MM-DD)")
	endStr := flag.String("end", "", "End date (YYYY-MM-DD)")
	flag.Parse()

	if *sysPath == "" || *bankPaths == "" || *startStr == "" || *endStr == "" {
		log.Fatal("Usage: go run . -system=sys.csv -banks=b1.csv,b2.csv -start=2025-01-01 -end=2025-01-31")
	}

	startDate, _ := time.Parse("2006-01-02", *startStr)
	endDate, _ := time.Parse("2006-01-02", *endStr)

	sysTrxs, err := ParseSystemCSV(*sysPath)
	if err != nil {
		log.Fatalf("System parse error: %v", err)
	}
	// Filter each bank file by date
	filteredSys := FilterSystemByDate(sysTrxs, startDate, endDate)

	var allBankRecs []BankRecord
	bankFileList := strings.Split(*bankPaths, ",")

	for _, path := range bankFileList {
		recs, err := ParseBankCSV(strings.TrimSpace(path))
		if err != nil {
			log.Printf("Warning: Could not read bank file %s: %v", path, err)
			continue
		}
		// Filter each bank file by date
		filteredBank := FilterBankByDate(recs, startDate, endDate)
		allBankRecs = append(allBankRecs, filteredBank...)
	}

	summary := Reconcile(filteredSys, allBankRecs, *windowDays)
	fmt.Printf("\n--- AMARTHA RECONCILIATION REPORT ---\n")
	fmt.Printf("Matched : %d\n", summary.MatchedCount)
	fmt.Printf("Unmatched : %d\n", summary.UnmatchedCount)
	fmt.Printf("Discrepancy Total: Rp.%d,-\n", summary.TotalDiscrepancy)
}
