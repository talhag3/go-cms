/* Go Concurrency: Timeout with Select */
package main

import (
	"fmt"
	"time"
)

func main() {
	// 1 CREATE A BUFFERED CHANNEL (capacity = 1)
	//    - Can hold 1 value without blocking
	//    - Like a mailbox that can hold 1 letter
	c1 := make(chan string, 1)

	// 2️ LAUNCH A GOROUTINE (background worker)
	//    - `go` keyword starts a NEW thread/goroutine
	//    - This runs IN PARALLEL with main()
	go func() {
		// Simulate slow work (2 seconds)
		time.Sleep(2 * time.Second)

		// Send data into the channel
		c1 <- "result 1"
	}()

	// 3️ SELECT STATEMENT (like a switch for channels)
	//    - Waits for MULTIPLE operations simultaneously
	//    - Executes the FIRST one that's ready
	select {
	case res := <-c1:
		// Case A: Data received from c1
		fmt.Println(res)

	case <-time.After(1 * time.Second):
		// Case B: TIMEOUT after 1 second
		//    - time.After() returns a channel
		//    - That channel sends current time after duration
		fmt.Println("timeout 1")
	}
}

/* Is time a channel? */

/*
Almost correct! Let me clarify:

time.After() returns a CHANNEL:
*/

/*

case <-time.After(1 * time.Second):
//       ^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//       This expression RETURNS a channel (<-chan Time)
//
// That channel will:
// - Block (wait) for 1 second
// - Then SEND the current time into itself
// - Your select receives from it → triggers timeout

*/

/*
// Simplified mental model of what time.After does:
func After(d Duration) <-chan Time {
    ch := make(chan Time, 1)  // Creates a channel

    go func() {                // Starts a timer goroutine
        time.Sleep(d)          // Waits for duration
        ch <- Now()            // Sends time into channel
    }()

    return ch                  // Returns the channel immediately
}

So you're receiving from a channel that time.After() created for you!


*/
