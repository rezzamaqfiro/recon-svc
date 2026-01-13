package main

import (
	"math"
	"strconv"
	"time"
)

func Reconcile(systemTrxs []Transaction, bankRecords []BankRecord, windowDays int) *ReconSummary {
	summary := &ReconSummary{BankUnmatched: make(map[string][]BankRecord)}

	// group system trx by date
	sysDateMap := make(map[string][]Transaction)
	for _, t := range systemTrxs {
		dateKey := t.Time.Format("2006-01-02")
		sysDateMap[dateKey] = append(sysDateMap[dateKey], t)
	}

	for _, b := range bankRecords {
		summary.TotalProcessed++

		var bestMatch *Transaction
		var bestDateKey string
		var bestIndex int
		var smallestDiff int64 = -1

		// search within the time window, constrain if windowDays is 1, it checks Date-1 to Date+1
		for i := -windowDays; i <= windowDays; i++ {
			searchDate := b.Date.AddDate(0, 0, i).Format("2006-01-02")
			candidates := sysDateMap[searchDate]

			for idx, sysTrx := range candidates {
				diff := int64(math.Abs(float64(sysTrx.Amount - b.Amount)))

				// find the candidate accross all window days with smallest diff
				if smallestDiff == -1 || diff < smallestDiff {
					smallestDiff = diff
					bestMatch = &candidates[idx]
					bestDateKey = searchDate
					bestIndex = idx
				}
			}
		}

		if bestMatch != nil {
			summary.MatchedCount++
			summary.TotalDiscrepancy += smallestDiff

			// Remove from map to prevent re-matching
			list := sysDateMap[bestDateKey]
			sysDateMap[bestDateKey] = append(list[:bestIndex], list[bestIndex+1:]...)
		} else {
			summary.UnmatchedCount++
			summary.BankUnmatched["Bank_A"] = append(summary.BankUnmatched["Bank_A"], b)
		}
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
