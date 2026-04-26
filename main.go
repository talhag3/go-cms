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

	go func(numsChan chan<- int) {
		for i := 1; i <= 5; i++ {
			numsChan <- i
		}
	}(nums)

	go func(numsChan <-chan int, squareChan chan<- int) {
		for value := range numsChan {
			squareChan <- (value * value)
		}
	}(nums, square)

	for i := 1; i <= 5; i++ {
		fmt.Println(<-square)
	}

}
