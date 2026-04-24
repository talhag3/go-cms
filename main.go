package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	// This loops until jobs channel is closed AND empty
	for job := range jobs {
		time.Sleep(1 * time.Second)
		fmt.Printf("Worker %d processing job %d\n", id, job)
		results <- job * 2 // Send result back (blocks if results is full)
	}
}

func main() {
	jobs := make(chan int, 5)    // Buffered pipe (can hold 5 items)
	results := make(chan int, 5) // Buffered pipe (can hold 5 items)

	// Start 3 workers (goroutines)
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send 5 jobs (fills the buffer)
	for j := 1; j <= 5; j++ {
		jobs <- j // Send to pipe (blocks if buffer full)
	}

	close(jobs) // Signal: "No more jobs coming"

	// Read 5 results (must match number of jobs sent)
	for a := 1; a <= 5; a++ {
		fmt.Println("Result:", <-results) // Blocks until result available
	}
}
