package main

import (
	"time"
)

func processInvestments(investments []Investment) (map[*Investment]Evaluation, error) {
	evaluations := make(map[*Investment]Evaluation)
	for i := range investments {
		inv := &investments[i]
		evaluation, err := processInvestment(inv)
		if err != nil {
			return nil, err
		}
		evaluations[inv] = evaluation
	}
	return evaluations, nil
}

func processInvestment(investment *Investment) (Evaluation, error) {
	startDate := investment.Transactions[0].Date
	evaluation := Evaluation{}
	evaluation.CurrentValue = investment.CurrentValue
	if len(investment.Transactions) == 0 {
		return evaluation, nil
	}
	durationDays, err := getDateDifferenceDays(startDate, "")
	if err != nil {
		return evaluation, err
	}
	evaluation.DurationDays = durationDays

	totalDeposits := 0.0
	totalWithdrawals := 0.0
	periodInvestment := 0.0
	totalDuration := 0
	timelyInvestment := 0.0

	for i, tx := range investment.Transactions {
		switch tx.Action {
		case DepositAction:
			totalDeposits += tx.Amount
			periodInvestment += tx.Amount
		case WithdrawAction:
			totalWithdrawals += tx.Amount
			periodInvestment -= tx.Amount
		case DividendAction, InterestAction:
			evaluation.Gain += tx.Amount
		}
		finalTransaction := (i == len(investment.Transactions)-1)
		if tx.Action == DepositAction || tx.Action == WithdrawAction || finalTransaction {
			var periodEndDate string
			if finalTransaction {
				periodEndDate = ""
			} else {
				periodEndDate = investment.Transactions[i+1].Date
			}
			period, err := getDateDifferenceDays(tx.Date, periodEndDate)
			if err != nil {
				return evaluation, err
			}
			totalDuration += period
			timelyInvestment += periodInvestment * float64(period)
		}
	}
	evaluation.NetInvested = totalDeposits - totalWithdrawals
	evaluation.TotalDeposits = totalDeposits
	evaluation.TotalWithdrawals = totalWithdrawals
	evaluation.Gain = evaluation.Gain + evaluation.CurrentValue - evaluation.NetInvested
	evaluation.GainPct = (evaluation.Gain / (timelyInvestment / float64(totalDuration))) * 100
	evaluation.GainAnnualizedPct = (evaluation.GainPct / float64(evaluation.DurationDays)) * 365

	return evaluation, nil
}

func getDateDifferenceDays(startDate, endDate string) (int, error) {
	d1, err := time.Parse(DateLayout, startDate)
	if err != nil {
		return 0, err
	}

	var d2 time.Time
	if endDate == "" {
		now := time.Now()
		d2 = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		d2, err = time.Parse(DateLayout, endDate)
		if err != nil {
			return 0, err
		}
	}

	days := int(d2.Sub(d1).Hours() / 24)

	if days < 1 {
		days = 1
	}
	return days, nil
}
