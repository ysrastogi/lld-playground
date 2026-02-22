package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Printf("Initial goroutines: %d\n", runtime.NumGoroutine())

	// Create a channel but we will never receive from it
	ch := make(chan int)

	// Launch multiple goroutines that will try to send to the channel
	for i := 0; i < 100; i++ {
		go func(id int) {
			// This will block forever because there is no receiver.
			// The goroutines are now leaked and will occupy memory indefinitely.
			ch <- id
		}(i)
	}

	// Give them a moment to start and get stuck
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("Goroutines after leak: %d\n", runtime.NumGoroutine())

	// In a real long-running application, these leaked goroutines would accumulate
	// and eventually cause an Out Of Memory (OOM) crash.
}
