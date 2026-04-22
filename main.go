package main

import (
	"fmt"
	"time"
)

func sayHello(name string) {
	time.Sleep(1 * time.Second)
	fmt.Println("Hello", name)
}

func main() {
	// Regular function call (blocking)
	sayHello("Ali") // Takes 1 second

	// Goroutine (non-blocking)
	go sayHello("Bilal") // Runs in background
	go sayHello("Saleem")

	time.Sleep(2 * time.Second) // Wait for goroutines to finish
	fmt.Println("Done")
}
