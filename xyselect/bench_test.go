package xyselect_test

import (
	"testing"

	"github.com/xybor/xyplatform/xyselect"
)

func runBenchSelector(b *testing.B, selector *xyselect.Selector) {
	var c = make(chan int)
	go func() {
		for i := 0; i < b.N; i++ {
			c <- i
		}
		close(c)
	}()

	selector.Recv(xyselect.C(c))
	for i := 0; i < b.N; i++ {
		selector.Select(false)
	}
}

func BenchmarkRSelector(b *testing.B) {
	runBenchSelector(b, xyselect.R())
}

func BenchmarkESelector(b *testing.B) {
	runBenchSelector(b, xyselect.E())
}
