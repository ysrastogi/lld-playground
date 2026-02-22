package main

import "testing"

// Run with: go test -bench=.

// --- Interface Setup ---
type Shape interface {
	Area() float64
}

// --- Concrete Type Setup ---
type Rectangle struct {
	width, height float64
}

//go:noinline
func (r *Rectangle) Area() float64 {
	return r.width * r.height
}

// Function accepting an interface (dynamic dispatch)
//
//go:noinline
func calculateAreaInterface(s Shape) float64 {
	return s.Area()
}

// Function accepting a concrete type (static dispatch)
//
//go:noinline
func calculateAreaConcrete(r *Rectangle) float64 {
	return r.Area()
}

func BenchmarkInterfaceCall(b *testing.B) {
	rect := &Rectangle{width: 10, height: 5}
	var shape Shape = rect

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calculateAreaInterface(shape)
	}
}

func BenchmarkConcreteCall(b *testing.B) {
	rect := &Rectangle{width: 10, height: 5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calculateAreaConcrete(rect)
	}
}
