/*
Keep receiving until channel closes

Create a goroutine that sends numbers 1 to 5, then closes the channel.
Use select inside a loop to receive values.

Hint: When receiving from a closed channel:

value, ok := <-ch

ok becomes false.
*/

package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 5)

	go func(_ch chan<- int) {

		for i := 1; i <= 5; i++ {
			time.Sleep(1 * time.Second)
			ch <- i
		}
		close(_ch)
	}(ch)

	for {
		select {
		case value, ok := <-ch:
			if !ok {
				fmt.Println("channel closed")
				break
			}
			fmt.Println(value)
		}
	}
}
