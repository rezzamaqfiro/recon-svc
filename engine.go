package main

import (
	"time"
)

// I use a map to Store system transactions for high speed lookups
func Reconcile(systemTrxs []Transaction, bankRecords []BankRecord) *ReconSummary {
	summary := &ReconSummary{
		SystemUnmatched: []Transaction{},
		BankUnmatched: make(map[string][]BankRecord),
	}

	// Step 1: Index system transactions by a "Recon Key"
	// Since BankIDs don't match system IDs, we match by Amount and Date
	// in a real scenario, you'd use a more complex key lor a sliding date window
	sysMap := make(map[string]Transaction)
	for _, t := range systemTrxs {
		key := generateKey(t.Amount, t.Time)
		sysMap[key] = t
		summary.TotalProcessed++
	}

	// Step 2: itterate through bank records to find matches
	for _, b := range bankRecords {
		key := generateKey(b.Amount, b.Date)

		if sysTrx, found := sysMap[key]; found {
			// SUCCESS: Match Found
			summary.MatchedCount++
			// Calcilate discrepancy if any (per requirements)
			diff := int64(math.Abs(float64(sysTrx.Amount - b.Amount)))
			summary.TotalDiscrepancy += diff

			// Remove from map so we know what's left is unmatched
			delete(sysMap, key)
		} else {
			// FAILURE: Bank Record no system equivalent
			summary.BankUnmatched["Bank_A"] = append(summary.BankUnmatched["Bank_A"], b)
			summary.UnmatchedCount++
		}

	}
	
}

// generateKey creates a unique string based on amount and date
func generateKey(amount int64, t time.Time) string {
	return t.Format("2006-01-02") + "_" + string(amount)
}