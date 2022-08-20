package xycond_test

import (
	"testing"

	"github.com/xybor/xyplatform/xycond"
)

var a []int

func init() {
	for i := 0; i < 100000; i++ {
		a = append(a, i)
	}
}

func BenchmarkMustContainA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xycond.MustContainA(a, i)
	}
}
