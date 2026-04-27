package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	// CONCEPT: WaitGroup
	// Think of this as a counter for background jobs.
	// It prevents the main program from exiting before background tasks finish.
	var wg sync.WaitGroup
	wg.Add(2) // "Hey Go, expect 2 background tasks to finish before you stop the program"

	// CONCEPT: Buffered Channels (Size 1)
	// These act like mailboxes that can hold exactly 1 letter.
	// Because they have a buffer, the goroutines can drop their message in the box
	// and immediately finish, without waiting for main() to read it.
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	// CONCEPT: Goroutine 1
	// The 'go' keyword launches this in the background (like a non-blocking async call in PHP)
	go func(ch chan<- string) {
		defer wg.Done() // PROMISE: No matter how this function exits, subtract 1 from the WaitGroup

		fmt.Println("go1 running")
		time.Sleep(1 * time.Second)    // Simulate a slow task (e.g., a 1-second database query)
		ch <- "Go routine 1 - sending" // Drop the result into the mailbox

		fmt.Println("go1 end")
	}(ch1)

	// CONCEPT: Goroutine 2
	// This runs simultaneously with Goroutine 1
	go func(ch chan<- string) {
		defer wg.Done() // PROMISE: Subtract 1 from the WaitGroup when done

		fmt.Println("go2 running")
		time.Sleep(2 * time.Second)    // Simulate an even SLOWER task (2 seconds)
		ch <- "Go routine 2 - sending" // Drop the result into the mailbox

		fmt.Println("go2 end")
	}(ch2)

	// CONCEPT: Select Statement (The "Race Receiver")
	// This is like JavaScript's Promise.race().
	// It pauses main() and listens to BOTH channels.
	// The millisecond a message arrives in EITHER channel, it runs that specific 'case',
	// grabs the message, and then EXITS the select block entirely.
	select {
	case msg := <-ch1:
		fmt.Println("Received:", msg)
	case msg := <-ch2:
		fmt.Println("Received:", msg)
	}

	// ⚠️ THE "GOTCHA" CONCEPT: What happens to the loser?
	// Because select is NOT a loop, it only runs ONCE.
	// Goroutine 1 finishes in 1 second. Goroutine 2 finishes in 2 seconds.
	// Select will grab ch1's message and leave.
	//
	// But what about ch2?
	// Goroutine 2 WILL finish after 2 seconds. It WILL put its message in ch2.
	// But NOBODY will ever read it. The message just sits in the buffer forever
	// until the program ends and the garbage collector cleans it up.

	// CONCEPT: Blocking Wait
	// Main() reaches here and freezes.
	// It will wait here until wg.Done() has been called exactly 2 times.
	// Even though select ignored go2's message, go2 itself still finishes its sleep
	// and calls wg.Done(), allowing the program to safely exit.
	wg.Wait()

	fmt.Println("all goroutines finished")
}
