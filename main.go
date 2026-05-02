/*
Write a function fetchData(ctx context.Context) string that simulates a slow database query by time.Sleep(3 * time.Second) and
then returns "Data fetched".
In main, call fetchData, but enforce a 1-second timeout. Print the result or the error.
*/

package main

import (
	"context"
	"fmt"
	"time"
)

func fetchData(ctx context.Context) (string, error) {
	resultChan := make(chan string)

	go func() {
		time.Sleep(3 * time.Second) // Slow query
		resultChan <- "Data fetched"
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err() // Returns "context deadline exceeded"
	case res := <-resultChan:
		return res, nil
	}
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	fetchData(ctx)

	result, err := fetchData(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(result)
	}
}
