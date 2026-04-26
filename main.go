/*
Square numbers

Create:

one goroutine that sends numbers 1 to 5
another goroutine that receives numbers and sends their squares
main prints the squared numbers

Example output:

1
4
9
16
25

Hint: Use two channels: nums and squares.

*/

package main

import "fmt"

func main() {

	nums := make(chan int)
	square := make(chan int)

	// Goroutine 1: The Producer
	go func(numsChan chan<- int) {
		for i := 1; i <= 5; i++ {
			numsChan <- i
		}
		close(numsChan) // FIX: Tell the receiver "I'm done sending"
	}(nums)

	// Goroutine 2: The Middleman
	go func(numsChan <-chan int, squareChan chan<- int) {
		for value := range numsChan { // ✅ Now this safely stops when nums is closed
			squareChan <- (value * value)
		}
		close(squareChan) // It's also good practice to close the output channel
	}(nums, square)

	// Main: The Consumer
	for i := 1; i <= 5; i++ {
		fmt.Println(<-square)
	}

}

/*

In Go, when the main() function finishes, the entire program exits instantly. It doesn't care if other goroutines are still running or waiting. It pulls the plug.

So, your second goroutine was about to freeze forever waiting for close(nums), but main() ended a millisecond before it could become a problem.

Why this is a dangerous bug
If you ever upgrade this code to use a sync.WaitGroup (a tool to make main() wait for all goroutines to finish), your program will deadlock and freeze. main() will wait for the second goroutine to finish, and the second goroutine will wait for close(nums). Forever.
*/
