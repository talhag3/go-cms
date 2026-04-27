/*
What causes blocking?

Create a buffered channel with size 2.

Try sending:

10
20
30

without receiving.

explain this problem to me with code
*/

package main

import "fmt"

func main() {
	// Create a box that only holds exactly 2 items
	ch := make(chan int, 2)

	ch <- 10 // OK
	ch <- 20 // OK

	// 🚨 THE PROBLEM IS HERE 🚨
	ch <- 30

	fmt.Println("I will never print this")
}

/*
The Golden Rule of Go: If you send to a buffered channel more times than its buffer size, you MUST have another goroutine receiving from it at the same time, or your program will deadlock.
*/
