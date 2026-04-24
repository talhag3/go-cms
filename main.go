package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 1; i < 6; i++ {
			fmt.Println("hello")
			time.Sleep(2 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()

		for i := 1; i < 6; i++ {
			fmt.Println("world")
			time.Sleep(2 * time.Millisecond)
		}
	}()

	wg.Wait()

	fmt.Println("all goroutines finished")
}
