package main

import (
	"log"
)

func main() {
	investments, err := loadCurrentInvestments()
	if err != nil {
		log.Fatalf("Failed to load current investments: %v", err)
	}
	for _, inv := range investments {
		log.Printf("%s: %d", inv.Name, len(inv.Transactions))
	}
}
