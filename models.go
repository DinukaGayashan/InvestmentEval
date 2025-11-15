package main

type Investment struct {
	Name         string
	SheetName    string
	CurrentValue float64
	Transactions []Transaction
}

type Transaction struct {
	Date   string
	Action string
	Amount float64
}

type Evaluation struct {
	CurrentValue      float64
	NetInvested       float64
	TotalDeposits     float64
	TotalWithdrawals  float64
	DurationDays      int
	Gain              float64
	GainPct           float64
	GainAnnualizedPct float64
}
