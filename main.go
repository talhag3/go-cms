/*
Create a context that holds a userID (string) and an agentLevel (int). Pass this context to a function called secretMission.
The function must extract both values and print: "Agent [userID] executing mission at level [agentLevel]"
*/

package main

import (
	"context"
	"fmt"
)

// Custom type prevents key collisions
type contextKey string

const (
	userIDKey     contextKey = "user_id"
	agentLevelKey contextKey = "agent_level"
)

func secretMission(ctx context.Context) {
	// Extract values safely
	userID, ok1 := ctx.Value(userIDKey).(string)
	level, ok2 := ctx.Value(agentLevelKey).(int)

	if !ok1 || !ok2 {
		fmt.Println("Mission aborted: Missing credentials")
		return
	}

	fmt.Printf("Agent %s executing mission at level %d\n", userID, level)
}

func main() {

	ctx := context.Background()

	// Chain values into context
	ctx = context.WithValue(ctx, userIDKey, "James Bond")
	ctx = context.WithValue(ctx, agentLevelKey, 7)

	secretMission(ctx)
	// time.Sleep(time.Second * 2)
}
