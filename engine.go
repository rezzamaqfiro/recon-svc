package main

import (
	"math"
	"strconv"
	"time"
)

// I use a map to Store system transactions for high speed lookups
func Reconcile(systemTrxs []Transaction, bankRecords []BankRecord) *ReconSummary {
	summary := &ReconSummary{BankUnmatched: make(map[string][]BankRecord)}

	sysDateMap := make(map[string][]Transaction)
	for _, t := range systemTrxs {
		dateKey := t.Time.Format("2006-01-02")
		sysDateMap[dateKey] = append(sysDateMap[dateKey], t)
	}

	for _, b := range bankRecords {
		summary.TotalProcessed++
		dateKey := b.Date.Format("2006-01-02")

		candidates, found := sysDateMap[dateKey]
		if !found || len(candidates) == 0 {
			summary.UnmatchedCount++
			summary.BankUnmatched["Bank_A"] = append(summary.BankUnmatched["Bank_A"], b)
			continue
		}

		// HEURISTIC
		bestIndex := -1
		var smallestDiff int64 = -1

		for i, sysTrx := range candidates {
			diff := int64(math.Abs(float64(sysTrx.Amount - b.Amount)))
			if smallestDiff == -1 || diff < smallestDiff {
				smallestDiff = diff
				bestIndex = i
			}
		}

		summary.MatchedCount++
		summary.TotalDiscrepancy += smallestDiff

		// Remove the matched system record so it's not used again
		sysDateMap[dateKey] = append(candidates[:bestIndex], candidates[bestIndex+1:]...)
	}

	for _, list := range sysDateMap {
		for _, t := range list {
			summary.UnmatchedCount++
			summary.SystemUnmatched = append(summary.SystemUnmatched, t)
		}
	}

	return summary
}

// generateKey creates a unique string based on amount and date
func generateKey(amount int64, t time.Time) string {
	return t.Format("2006-01-02") + "_" + strconv.FormatInt(amount, 10)
}
