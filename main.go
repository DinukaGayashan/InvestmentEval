package main

import (
	"context"
	"log"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	spreadsheet, err := srv.Spreadsheets.Get(SpreadsheetID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve spreadsheet: %v", err)
	}

	for _, sheet := range spreadsheet.Sheets {
		if strings.HasPrefix(sheet.Properties.Title, CurrentPrefix) {
			log.Println(sheet.Properties.Title)
		}
	}
}
