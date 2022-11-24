package recursion

import (
	"fmt"
	"testing"
)

func TestFib(t *testing.T) {
	fn := Fib()
	for i := 0; i < 10; i++ {
		fmt.Println(fn())
	}
}

func TestFibIter(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(FibIter(i + 1))
	}
}

func TestFibRec(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(FibRec(i + 1))
	}
}

func TestFibTailRec(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(FibTailRec(i+1, 0, 1))
	}
}

// Result average: 102m iterations, 11.4 ns/op, 0 B/op, 0 allocs/op
func BenchmarkFib(b *testing.B) {
	b.ReportAllocs()
	var res, save int
	for i := 0; i < b.N; i++ {
		fn := Fib()
		for j := 0; j < 10; j++ {
			res = fn()
			if res < 0 {
				b.Errorf("something went wrong")
			}
		}
		save = res
	}
	_ = save
	_ = res
}

// Result average: 29m iterations, 43.4 ns/op, 0 B/op, 0 allocs/op
func BenchmarkFibIter(b *testing.B) {
	b.ReportAllocs()
	var res, save int
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			res = FibIter(j + 1)
			if res < 0 {
				b.Errorf("something went wrong")
			}
		}
		save = res
	}
	_ = save
	_ = res
}

// Result average: 1.3m iterations, 778.4 ns/op, 0 B/op, 0 allocs/op
func BenchmarkFibRec(b *testing.B) {
	b.ReportAllocs()
	var res, save int
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			res = FibRec(j + 1)
			if res < 0 {
				b.Errorf("something went wrong")
			}
		}
		save = res
	}
	_ = save
	_ = res
}

// Result average: 12m iterations, 86.65 ns/op, 0 B/op, 0 allocs/op
func BenchmarkFibTailRec(b *testing.B) {
	b.ReportAllocs()
	var res, save int
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			res = FibTailRec(j+1, 0, 1)
			if res < 0 {
				b.Errorf("something went wrong")
			}
		}
		save = res
	}
	_ = save
	_ = res
}
