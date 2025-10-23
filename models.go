package main

type Investment struct {
	Name         string
	SheetName    string
	Transactions []Transaction
}

type Transaction struct {
	Date   string
	Action string
	Amount float64
}

type Evaluation struct {
	CurrentValue      float64
	DurationDays      int
	Gain              float64
	GainPct           float64
	GainAnnualizedPct float64
}
