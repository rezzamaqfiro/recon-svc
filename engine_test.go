package main

import (
	"testing"
	"time"
)

func TestReconcile(t *testing.T) {
	tests := []struct {
		name              string
		systemData        []Transaction
		bankData          []BankRecord
		windowDays        int
		expectedMatched   int
		expectedUnmatched int
	}{
		{
			name: "Perfect Match",
			systemData: []Transaction{
				{TrxID: "1", Amount: 10000, Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			bankData: []BankRecord{
				{ID: "A", Amount: 10000, Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			windowDays:        0,
			expectedMatched:   1,
			expectedUnmatched: 0,
		},
		{
			name: "Unmatched System Record",
			systemData: []Transaction{
				{TrxID: "2", Amount: 50000, Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			bankData:          []BankRecord{},
			windowDays:        0,
			expectedMatched:   0,
			expectedUnmatched: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := Reconcile(tt.systemData, tt.bankData, tt.windowDays)
			if summary.MatchedCount != tt.expectedMatched {
				t.Errorf("expected %d matches, got %d", tt.expectedMatched, summary.MatchedCount)
			}
			if summary.UnmatchedCount != tt.expectedUnmatched {
				t.Errorf("expected %d matches, got %d", tt.expectedUnmatched, summary.UnmatchedCount)
			}
		})
	}
}

func TestReconcile_Advanced(t *testing.T) {
	// set the dates
	day1 := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	day2 := time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name                string
		systemData          []Transaction
		bankData            []BankRecord
		windowDays          int
		expectedMatched     int
		expectedUnmatched   int
		expectedDiscrepancy int64
	}{
		{
			name: "Scenario 1: Exact Match",
			systemData: []Transaction{
				{TrxID: "ST1", Amount: 20000, Time: day1},
			},
			bankData: []BankRecord{
				{ID: "SB1", Amount: 20000, Date: day1},
			},
			windowDays:          0,
			expectedMatched:     1,
			expectedUnmatched:   0,
			expectedDiscrepancy: 0,
		},
		{
			name: "Scenario 2: Amount Discrepancy",
			systemData: []Transaction{
				{TrxID: "ST2", Amount: 125000, Time: day1},
			},
			bankData: []BankRecord{
				{ID: "SB2", Amount: 124400, Date: day1},
			},
			windowDays:          0,
			expectedMatched:     1,
			expectedUnmatched:   0,
			expectedDiscrepancy: 600, // 125000 - 124400
		},
		{
			name: "Scenario 3: Time Window Match (Settlement Lag)",
			systemData: []Transaction{
				{TrxID: "ST3", Amount: 75000, Time: day1},
			},
			bankData: []BankRecord{
				{ID: "SB3", Amount: 75000, Date: day2}, // Next day
			},
			windowDays:          1,
			expectedMatched:     1,
			expectedUnmatched:   0,
			expectedDiscrepancy: 0,
		},
		{
			name: "Scenario 4: Total Unmatched Records",
			systemData: []Transaction{
				{TrxID: "ST4", Amount: 50000, Time: day1},
			},
			bankData: []BankRecord{
				{ID: "SB4", Amount: 88888, Date: day2},
			},
			windowDays:          0,
			expectedMatched:     0,
			expectedUnmatched:   2, // both records unmatched
			expectedDiscrepancy: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := Reconcile(tt.systemData, tt.bankData, tt.windowDays)

			if summary.MatchedCount != tt.expectedMatched {
				t.Errorf("MatchedCount: got %d, want %d", summary.MatchedCount, tt.expectedMatched)
			}
			if summary.UnmatchedCount != tt.expectedUnmatched {
				t.Errorf("UnmatchedCount: got %d, want %d", summary.UnmatchedCount, tt.expectedUnmatched)
			}
			if summary.TotalDiscrepancy != tt.expectedDiscrepancy {
				t.Errorf("Discrepancy: got %d, want %d", summary.TotalDiscrepancy, tt.expectedDiscrepancy)
			}
		})
	}
}
