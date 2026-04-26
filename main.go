// Practice  Send multiple values

package main

import "fmt"

func main() {

	myChan := make(chan int)

	go func(dataChan chan<- int) {

		for i := 1; i < 6; i++ {
			dataChan <- i
		}

	}(myChan)

	for i := 1; i < 6; i++ {
		fmt.Println(<-myChan)
	}

	fmt.Println("End Program")
}
