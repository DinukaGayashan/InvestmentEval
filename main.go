package main

import (
	"log"
)

func main() {
	investments, err := loadCurrentInvestments()
	if err != nil {
		log.Fatalf("Failed to load current investments: %v", err)
	}
	statistics, err := processInvestments(investments)
	if err != nil {
		log.Fatalf("Failed to process investments: %v", err)
	}
	err = uploadStatistics(statistics)
	if err != nil {
		log.Fatalf("Failed to upload statistics: %v", err)
	}
}
