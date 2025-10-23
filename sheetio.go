package main

import (
	"context"
	"strings"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	spreadsheet *sheets.Spreadsheet
	loadOnce    sync.Once
	loadError   error
)

func GetSpreadsheet() (*sheets.Spreadsheet, error) {
	loadOnce.Do(func() {
		ctx := context.Background()
		sheetSrv, err := sheets.NewService(ctx, option.WithCredentialsFile(CredentialsFile))
		if err != nil {
			loadError = err
			return
		}
		spreadsheet, loadError = sheetSrv.Spreadsheets.Get(SpreadsheetID).Do()
	})
	if loadError != nil {
		return nil, loadError
	}
	return spreadsheet, nil
}

func loadCurrentInvestments() ([]Investment, error) {
	investments := []Investment{}
	spreadsheet, err := GetSpreadsheet()
	if err != nil {
		return nil, err
	}
	for _, sheet := range spreadsheet.Sheets {
		if strings.HasPrefix(sheet.Properties.Title, CurrentPrefix) {
			investments = append(investments, Investment{
				Name:      strings.TrimPrefix(sheet.Properties.Title, CurrentPrefix),
				SheetName: sheet.Properties.Title,
			})
		}
	}
	return investments, nil
}
