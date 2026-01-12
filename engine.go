package main

import (
	"time"
)

// generateKey creates a unique string based on amount and date
func generateKey(amount int64, t time.Time) string {
	return t.Format("2006-01-02") + "_" + string(amount)
}