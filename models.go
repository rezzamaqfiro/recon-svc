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

