package main

import (
	"fmt"
	"time"
)

func main() {
	// ============================================
	// 1️ CREATE A BUFFERED CHANNEL (capacity = 5)
	//    - Can hold up to 5 integers without blocking sender
	//    - Buffered = sender can send N values before blocking
	// ============================================
	ch := make(chan int, 5)

	// ============================================
	// 2️ LAUNCH PRODUCER GOROUTINE
	//    - Sends numbers 1-5 with 1-second delays
	//    - Closes channel when done (signals EOF)
	// ============================================
	go func(_ch chan<- int) {
		// Send numbers 1 through 5
		for i := 1; i <= 5; i++ {
			time.Sleep(1 * time.Second) // Simulate work/delay
			ch <- i                     // Send value into channel
			fmt.Printf("[Sender] Sent: %d\n", i)
		}

		// 🔚 CLOSE THE CHANNEL (important!)
		// - Signals receivers: "no more data"
		// - Future receives return (value, false)
		// - Channel cannot be re-used after closing
		close(_ch)
		fmt.Println("[Sender] Channel closed")
	}(ch) // Pass ch as argument to goroutine

	// ============================================
	// 3️ CONSUMER LOOP: Receive Until Channel Closes
	//    Uses labeled break to exit outer for-loop
	// ============================================

outerLoop: // ← LABEL for the for loop (needed to break out of both select AND for)
	for {
		select {
		case value, ok := <-ch:
			//  TWO-VALUE RECEIVE: value, ok := <-ch
			//    - value: the received data (or zero-value if closed)
			//    - ok: true if channel open, false if closed

			if !ok {
				// Channel is CLOSED!
				// - No more data will come
				// - Time to stop receiving

				fmt.Println("\n[Receiver] Channel detected as closed!")
				fmt.Println("[Receiver] Exiting receiver loop...")

				break outerLoop // ← LABELED BREAK! Exits the FOR loop entirely
				//   Without label: would only exit select (causing your bug!)
			}

			// Channel is OPEN and we got data
			fmt.Printf("[Receiver] Received: %d\n", value)

			// Note: No timeout case here, so we block waiting for data/close
		}
	}

	// ============================================
	// 4️ AFTER LOOP: Program continues here
	// ============================================
	fmt.Println("\n Program completed successfully!")
	fmt.Println(" All data received, channel was properly closed")
}
