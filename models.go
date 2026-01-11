package main

import "time"

type Transaction struct {
	TrxID 	string  
	Amount	int64
	Type	string	// "credit" or "debit"
	Time	time.Time
}

type BankRecord struct {
	ID		string
	Amount	int64
	Date	time.Time
}

type ReconSummary struct {
	TotalProcessed		int
	MatchedCount		int
	UnmatchedCount		int
	TotalDiscrepancy	int64			// Sum of absolute differences
	SystemUnmatched		[]Transaction	// Missing in bank records
	BankUnmatched		map[string][]BankRecord	// Missing in system records
}

