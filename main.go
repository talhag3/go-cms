package main

import (
	"fmt"
	"time"
)

// This is like a Laravel Job class
func processOrder(orderID int, result chan<- string) {
	time.Sleep(1 * time.Second) // Simulate work
	result <- fmt.Sprintf("Order %d processed!", orderID)
}

func main() {
	// Channel for results
	results := make(chan string, 3)

	// Dispatch 3 jobs (like dispatching Laravel jobs)
	go processOrder(101, results)
	go processOrder(102, results)
	go processOrder(103, results)

	// Wait for all results
	for i := 0; i < 3; i++ {
		fmt.Println(<-results)
	}
}
