package main

import (
	"testing"
)

func BenchmarkDay08(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}
