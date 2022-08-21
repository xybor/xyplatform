package xyselect_test

import (
	"testing"

	"github.com/xybor/xyplatform/xyselect"
)

func BenchmarkRSelector(b *testing.B) {
	var selector = xyselect.R()
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

func BenchmarkESelector(b *testing.B) {
	var selector = xyselect.E()
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
