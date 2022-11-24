package recursion

// Fib returns a function that returns
// successive Fibonacci numbers.
func Fib() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func FibIter(n int) int {
	if n <= 1 {
		return n
	}
	prev, current := 0, 1
	for i := 2; i <= n; i++ {
		prev, current = current, current+prev
	}
	return current
}

func FibRec(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return FibRec(n-1) + FibRec(n-2)
}

func FibTailRec(n, prev, current int) int {
	if n == 0 {
		return prev
	}
	if n == 1 {
		return current
	}
	return FibTailRec(n-1, current, prev+current)
}
