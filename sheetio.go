package main

import (
	"context"
	"strings"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	cachedSpreadsheet *sheets.Spreadsheet
	loadOnce          sync.Once
	loadError         error
)

func getSpreadsheet() (*sheets.Spreadsheet, error) {
	loadOnce.Do(func() {
		ctx := context.Background()
		sheetSrv, err := sheets.NewService(ctx, option.WithCredentialsFile(CredentialsFile))
		if err != nil {
			loadError = err
			return
		}
		cachedSpreadsheet, loadError = sheetSrv.Spreadsheets.Get(SpreadsheetID).Do()
	})
	if loadError != nil {
		return nil, loadError
	}
	return cachedSpreadsheet, nil
}

func loadCurrentInvestments() ([]Investment, error) {
	spreadsheet, err := getSpreadsheet()
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
		}
	}
	return investments, nil
}
