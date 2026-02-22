package main

import "fmt"

func main() {
	fmt.Println("Main started")

	// Create an unbuffered channel
	ch := make(chan int)

	// We are attempting to send to an unbuffered channel in the same goroutine.
	// This requires another goroutine to be ready to receive it *at the same time*.
	// Since there is no other goroutine, this single goroutine deadlocks instantly.
	ch <- 1

	fmt.Println("This line will never be reached")
}
