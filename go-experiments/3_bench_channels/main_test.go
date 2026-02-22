package main

import "testing"

// Run with: go test -bench=.

func BenchmarkUnbufferedChannel(b *testing.B) {
	ch := make(chan struct{})

	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
	}
}

func BenchmarkBufferedChannel(b *testing.B) {
	// A buffer size of 100
	ch := make(chan struct{}, 100)

	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
	}
}
