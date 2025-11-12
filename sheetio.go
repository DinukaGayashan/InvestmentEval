package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	sheetSrv *sheets.Service
	loadOnce sync.Once
	srvErr   error
)

func getService() (*sheets.Service, error) {
	loadOnce.Do(func() {
		ctx := context.Background()
		srv, err := sheets.NewService(ctx, option.WithCredentialsFile(CredentialsFile))
		if err != nil {
			srvErr = fmt.Errorf("failed to create sheets service: %v", err)
			return
		}
		sheetSrv = srv
	})
	return sheetSrv, srvErr
}

func loadCurrentInvestments() ([]Investment, error) {
	srv, err := getService()
	if err != nil {
		return nil, err
	}
	spreadsheet, err := srv.Spreadsheets.Get(SpreadsheetID).Do()
	if err != nil {
		return nil, err
	}

	investments := []Investment{}
	for _, sheet := range spreadsheet.Sheets {
		if after, ok := strings.CutPrefix(sheet.Properties.Title, CurrentPrefix); ok {
			investments = append(investments, Investment{
				Name:      after,
				SheetName: sheet.Properties.Title,
			})
			transactions, err := readTransactions(sheet.Properties.Title)
			if err != nil {
				return nil, err
			}
			investments[len(investments)-1].Transactions = transactions
		}
	}
	return investments, nil
}

func readTransactions(sheetName string) ([]Transaction, error) {
	srv, err := getService()
	if err != nil {
		return nil, err
	}

	readRange := sheetName + DataRange
	resp, err := srv.Spreadsheets.Values.Get(SpreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet %s: %v", sheetName, err)
	}

	transactions := []Transaction{}
	for i, row := range resp.Values {
		if len(row) < 3 {
			fmt.Printf("Missing values on record: %s:%v\n", sheetName, i)
			continue
		}
		dateStr := strings.TrimSpace(row[0].(string))
		actionStr := strings.TrimSpace(row[1].(string))
		amountStr := strings.TrimSpace(row[2].(string))

		amountStr = strings.ReplaceAll(amountStr, ",", "")
		amountVal := 0.0
		if val, err := strconv.ParseFloat(amountStr, 64); err == nil {
			amountVal = val
		} else {
			fmt.Println("Error parsing amount: ", err)
			continue
		}

		transaction := Transaction{
			Date:   dateStr,
			Action: actionStr,
			Amount: amountVal,
		}

		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
