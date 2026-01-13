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
	// Filter each system file by date
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
	fmt.Println("==========================================")
	fmt.Println("       AMARTHA RECONCILIATION REPORT       ")
	fmt.Println("==========================================")
	fmt.Printf("Total Processed:     %d\n", summary.TotalProcessed)
	fmt.Printf("Matched:             %d\n", summary.MatchedCount)
	fmt.Printf("Unmatched:           %d\n", summary.UnmatchedCount)
	fmt.Printf("Total Discrepancy:   Rp %d\n", summary.TotalDiscrepancy)
	fmt.Println("------------------------------------------")

	fmt.Println("\n MISSING IN BANK (Internal Records)")
	if len(summary.SystemUnmatched) == 0 {
		fmt.Println("      None - All system records found in bank.")
	} else {
		for i, trx := range summary.SystemUnmatched {
			if i >= 25 {
				fmt.Printf("      ... and %d more items\n", len(summary.SystemUnmatched)-25)
				break
			}
			fmt.Printf("      - ID: %-10s | Amt: %-10d | Date: %s\n",
				trx.TrxID, trx.Amount, trx.Time.Format("2006-01-02"))
		}
	}

	fmt.Println("\n MISSING IN SYSTEM (Bank Records)")
	if len(summary.BankUnmatched) == 0 {
		fmt.Println("      None - All bank records found in system.")
	} else {
		for bankName, records := range summary.BankUnmatched {
			fmt.Printf("      Bank: %s\n", bankName)
			for i, rec := range records {
				if i >= 25 {
					fmt.Printf("            ... and %d more items\n", len(records)-25)
					break
				}
				fmt.Printf("            - Ref: %-10s | Amt: %-10d | Date: %s\n",
					rec.ID, rec.Amount, rec.Date.Format("2006-01-02"))
			}
		}
	}
	fmt.Println("==========================================")
}
