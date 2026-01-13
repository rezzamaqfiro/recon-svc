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
			expectedMatched:   1,
			expectedUnmatched: 0,
		},
		{
			name: "Unmatched System Record",
			systemData: []Transaction{
				{TrxID: "2", Amount: 50000, Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			bankData:          []BankRecord{},
			expectedMatched:   0,
			expectedUnmatched: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := Reconcile(tt.systemData, tt.bankData)
			if summary.MatchedCount != tt.expectedMatched {
				t.Errorf("expected %d matches, got %d", tt.expectedMatched, summary.MatchedCount)
			}
			if summary.UnmatchedCount != tt.expectedUnmatched {
				t.Errorf("expected %d matches, got %d", tt.expectedUnmatched, summary.UnmatchedCount)
			}
		})
	}
}
