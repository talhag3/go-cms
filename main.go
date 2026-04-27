/*
Receive from two channels

Create two goroutines:

one sends "from channel 1"
one sends "from channel 2"

Use select in main to receive from whichever sends first.
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(2)

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	go func(ch chan<- string) {
		defer wg.Done()

		fmt.Println("go1 running")
		time.Sleep(1 * time.Second)
		ch <- "Go routine 1 - sending"

		fmt.Println("go1 end")
	}(ch1)

	go func(ch chan<- string) {
		defer wg.Done()

		fmt.Println("go2 running")
		time.Sleep(2 * time.Second)
		ch <- "Go routine 2 - sending"

		fmt.Println("go2 end")
	}(ch2)

	select {
	case msg := <-ch1:
		fmt.Println(msg)
	case msg := <-ch2:
		fmt.Println(msg)
	}

	wg.Wait()

	fmt.Println("all goroutines finished")
}
