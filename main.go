// Send one value through channel

package main

import (
	"fmt"
)

func main() {
	messageCh := make(chan string)

	go func(message chan<- string) {
		message <- "Hi Talha"
	}(messageCh)

	fmt.Println(<-messageCh)
}
