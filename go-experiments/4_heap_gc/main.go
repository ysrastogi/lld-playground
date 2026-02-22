package main

import (
	"fmt"
	"runtime"
)

// allocateOnHeap creates a large array and returns a pointer to it.
// Because the pointer escapes the function's scope, the Go compiler
// is forced to allocate this on the heap rather than the stack.
// Run with "go build -gcflags='-m'" to observe the escape analysis.
//
//go:noinline
func allocateOnHeap() *[1024 * 1024]int {
	arr := new([1024 * 1024]int)
	return arr
}

func main() {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)
	fmt.Printf("Before alloc - Heap Alloc: %v MB\n", m.Alloc/1024/1024)

	// We hold references to the allocated arrays so they aren't prematurely garbage collected.
	var hold [][1024 * 1024]int
	for i := 0; i < 20; i++ {
		hold = append(hold, *allocateOnHeap())
	}

	runtime.ReadMemStats(&m)
	fmt.Printf("After alloc  - Heap Alloc: %v MB\n", m.Alloc/1024/1024)

	// Remove references to make the allocated memory unreachable
	hold = nil

	// Manually trigger Garbage Collection
	runtime.GC()

	runtime.ReadMemStats(&m)
	fmt.Printf("After GC     - Heap Alloc: %v MB\n", m.Alloc/1024/1024)
}
