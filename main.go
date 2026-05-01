/* Create a goroutine that prints "working..."
every 500 milliseconds in an infinite loop. In your main function,
 let it run for exactly 2 seconds, then stop it.
 The program should exit cleanly without printing anything after it stops.
*/

package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker stopped!")
				return
			case <-time.After(500 * time.Millisecond):
				fmt.Println("working...")
			}
		}
	}(ctx)

	time.Sleep(time.Second * 2)

	// Send the stop signal
	cancel()

	// Give goroutine a moment to print "Worker stopped!"
	time.Sleep(100 * time.Millisecond)
}
