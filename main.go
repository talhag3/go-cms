/*

Buffered channel test

Create a buffered channel with size 2

*/

package main

import "fmt"

func main() {

	numBufChan := make(chan int, 2)

	go func(numChan chan<- int) {
		numChan <- 1
		numChan <- 2
		close(numChan)
	}(numBufChan)

	for v := range numBufChan {
		fmt.Println(v)
	}
}
