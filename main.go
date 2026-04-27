/*
Use timeout with select

Create a channel.
Start a goroutine that sends data after 3 seconds.
In main, use select with a 1-second timeout.

Expected idea:

timeout

Hint: Use time.After
*/

package main

import (
	"fmt"
	"time"
)

func main() {

	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout 1")
	}
}
