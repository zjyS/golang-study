package main

import (
	"testing"
)

func TestIsPrime(t *testing.T) {
	if !isPrime(13) {
		t.Error("13 is not prime")
	}
}

func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isPrime(9897)
	}
}
