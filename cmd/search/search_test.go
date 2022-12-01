package search

import (
	"testing"
	"time"
)

func TimeThis(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

const size = 64

var IntSlice = make(Slice[int], size, size)

func Benchmark_Slice_LinearSearch(b *testing.B) {
	for j := 0; j < len(IntSlice); j++ {
		IntSlice[j] = j
	}
	for i := 0; i < b.N; i++ {
		at, found := IntSlice.FindLinear(size - 8)
		if !found || IntSlice[at] != size-8 {
			b.Errorf("got=%d, wanted=%d\n", IntSlice[at], size-8)
		}
	}

}

func Benchmark_Slice_BinarySearch(b *testing.B) {
	for j := 0; j < len(IntSlice); j++ {
		IntSlice[j] = j
	}
	IntSlice.Sort()
	for i := 0; i < b.N; i++ {
		at, found := IntSlice.FindBinary(size - 8)
		if !found || IntSlice[at] != size-8 {
			b.Errorf("got=%d, wanted=%d\n", IntSlice[at], size-8)
		}
	}
}
