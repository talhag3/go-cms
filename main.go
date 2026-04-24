package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Response struct is used to communicate values and errors
// from a goroutine back to the main routine via channels [4].
type Response struct {
	value int
	err   error
}

func main() {
	start := time.Now() // Measure the total time the operation takes [5].

	// 1. Context with Value: Useful for request IDs or sharing state between goroutines [3, 6].
	ctx := context.WithValue(context.Background(), "foo", "bar")

	// 2. Context with Timeout: Ensures the function is deterministic and doesn't
	// hang longer than 200ms, regardless of how slow the third-party call is [2, 7].
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)

	// Always defer cancel to avoid context leaks [7].
	defer cancel()

	userID := 10
	val, err := fetchUserData(ctx, userID)
	if err != nil {
		log.Fatal(err) // If the context times out, we catch the error here [5, 8].
	}

	fmt.Printf("result: %d\n", val)
	fmt.Printf("took: %v\n", time.Since(start)) // Prints how long the execution took [5].
}

// fetchUserData mimics a business logic function that handles a request [9].
func fetchUserData(ctx context.Context, userID int) (int, error) {
	// Retrieve a value from the context (e.g., for tracing or logging) [3, 10].
	valFromCtx := ctx.Value("foo")
	fmt.Printf("Value from context (key 'foo'): %v\n", valFromCtx)

	// A channel is used to synchronize the result from the goroutine [4].
	respCh := make(chan Response)

	// We run the slow third-party call in a separate goroutine so we
	// can listen for the context timeout simultaneously [7].
	go func() {
		val, err := fetchThirdPartyStuffWhichCanBeSlow()
		respCh <- Response{
			value: val,
			err:   err,
		}
	}()

	// The select statement allows us to wait for multiple channel operations [11].
	for {
		select {
		case <-ctx.Done():
			// This case triggers if the context's timeout (200ms) is reached first [11, 12].
			return 0, fmt.Errorf("fetching data from third party took too long")
		case resp := <-respCh:
			// This case triggers if the third-party function returns before the timeout [13].
			return resp.value, resp.err
		}
	}
}

// fetchThirdPartyStuffWhichCanBeSlow mimics an external API call [1, 14].
func fetchThirdPartyStuffWhichCanBeSlow() (int, error) {
	// Replicate a delay (e.g., 150ms for success or 500ms to trigger a timeout) [6, 8, 9].
	time.Sleep(150 * time.Millisecond)

	return 666, nil // Returns a value (mimicking data) and no error [2, 9].
}
