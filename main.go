package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// requestIDMiddleware is a middleware function that adds a request ID to the context.
// This is useful for tracking requests across multiple handlers or services.
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new context with a request ID value.
		// This value can be accessed by any handler downstream.
		ctx := context.WithValue(r.Context(), "requestID", "12345")
		// Call the next handler in the chain, passing the new context.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// slowHandler simulates a slow database query.
// It demonstrates how to use context for cancellation.
func slowHandler(w http.ResponseWriter, r *http.Request) {
	// Get the context from the request.
	ctx := r.Context()

	// Simulate a slow operation (e.g., a database query) that takes 5 seconds.
	select {
	case <-time.After(5 * time.Second):
		// If the operation completes, retrieve the request ID from the context
		// and send a success response.
		requestID := ctx.Value("requestID").(string)
		fmt.Fprintf(w, "Slow query completed! Request ID: %s\n", requestID)
	case <-ctx.Done():
		// If the context is canceled (e.g., client disconnects or timeout),
		// log the error and send a 500 response.
		err := ctx.Err()
		log.Printf("Request canceled: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// timeoutHandler demonstrates how to use context for timeouts.
// It sets a 2-second timeout for the operation.
func timeoutHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new context with a 2-second timeout.
	// This context will automatically cancel after 2 seconds.
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	// Always call cancel() to release resources, even if the context times out.
	defer cancel()

	// Simulate an operation that takes 3 seconds.
	select {
	case <-time.After(3 * time.Second):
		// If the operation completes within 3 seconds, send a success response.
		fmt.Fprintf(w, "Timeout handler completed!\n")
	case <-ctx.Done():
		// If the context times out (after 2 seconds), log the error
		// and send a 408 (Request Timeout) response.
		err := ctx.Err()
		log.Printf("Timeout exceeded: %v\n", err)
		http.Error(w, err.Error(), http.StatusRequestTimeout)
	}
}

func main() {
	// Create a new HTTP request multiplexer.
	mux := http.NewServeMux()

	// Register the /slow endpoint with the requestIDMiddleware.
	// This ensures the request ID is added to the context before slowHandler is called.
	mux.Handle("/slow", requestIDMiddleware(http.HandlerFunc(slowHandler)))

	// Register the /timeout endpoint directly.
	mux.Handle("/timeout", http.HandlerFunc(timeoutHandler))

	// Start the HTTP server on port 8080.
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
