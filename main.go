package main

import (
	"log"
)

func main() {
	investments, err := loadCurrentInvestments()
	if err != nil {
		log.Fatalf("Failed to load current investments: %v", err)
	}
	for i, inv := range investments {
		log.Printf("Loading investment %d: %s", i+1, inv.Name)
	}
}
