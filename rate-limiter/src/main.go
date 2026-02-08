package main

import (
	"fmt"
	"rate-limiter/src/interfaces"
	"rate-limiter/src/models"
	"rate-limiter/src/services"
	"time"
)

func main() {
	stateStore := services.NewStateStore()
	policy := models.LimitPolicy{}
	policy.SetRequests(5)
	policy.SetTimeframe(10 * time.Second)
	policy.SetMaxBurst(5)
	policy.SetEntity(models.User)

	rule := &interfaces.TokenBucketRule{
		LimitPolicy: policy,
	}
	orchestrator := services.NewRateLimiterOrchestrator(stateStore, rule, policy)

	user1 := models.RequestContext{}
	user1.SetUserID("user-123")
	user1.SetAPIKey("api-key-abc")
	user1.SetIPAddress("192.168.1.1")

	user2 := models.RequestContext{}
	user2.SetUserID("user-456")
	user2.SetAPIKey("api-key-xyz")
	user2.SetIPAddress("192.168.1.2")

	fmt.Println("=== Rate Limiter Simulation ===")
	fmt.Printf("Policy: %d requests per %v\n\n", 5, 10*time.Second)

	fmt.Println("User 1 (user-123) making requests:")
	for i := 1; i <= 7; i++ {
		allowed := orchestrator.Allow(user1)
		status := "✓ ALLOWED"
		if !allowed {
			status = "✗ BLOCKED"
		}
		fmt.Printf("  Request %d: %s\n", i, status)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println()

	fmt.Println("User 2 (user-456) making requests:")
	for i := 1; i <= 4; i++ {
		allowed := orchestrator.Allow(user2)
		status := "✓ ALLOWED"
		if !allowed {
			status = "✗ BLOCKED"
		}
		fmt.Printf("  Request %d: %s\n", i, status)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println()

	fmt.Println("Waiting 11 seconds for token bucket to refill...")
	time.Sleep(11 * time.Second)
	fmt.Println("User 1 (user-123) making requests after refill:")
	for i := 1; i <= 3; i++ {
		allowed := orchestrator.Allow(user1)
		status := "✓ ALLOWED"
		if !allowed {
			status = "✗ BLOCKED"
		}
		fmt.Printf("  Request %d: %s\n", i, status)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n=== Simulation Complete ===")
}
