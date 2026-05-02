/*
A flash sale ends at exactly time.Now().Add(2 * time.Second).
Create a worker that checks the context's deadline, prints how many milliseconds are left, and loops every 500ms.
It should stop exactly when the sale ends.
*/

package main

import (
	"context"
	"fmt"
	"time"
)

func FlashSale(ctx context.Context) {
	for {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Sale is OVER!")
				return
			case <-time.After(500 * time.Millisecond):
				deadline, _ := ctx.Deadline()
				remaining := time.Until(deadline)
				fmt.Printf("Sale active... %v remaining\n", remaining.Round(time.Millisecond))
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancel()

	go FlashSale(ctx)

	// Wait longer than the sale to see it stop
	time.Sleep(3 * time.Second)
}
