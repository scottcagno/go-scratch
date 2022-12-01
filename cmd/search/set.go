package search

import (
	"sort"
)

type Comparable interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Slice[T Comparable] []T

func NewSlice[T Comparable](args ...int) Slice[T] {
	var length, capacity int
	if len(args) == 1 {
		length = args[0]
	}
	if len(args) == 2 {
		length = args[0]
		capacity = args[1]
	}
	return make(Slice[T], length, capacity)
}

func (s Slice[T]) Len() int {
	return len(s)
}

func (s Slice[T]) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Slice[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Slice[T]) Compare(this T, at int) int {
	if this == s[at] {
		return 0
	}
	if this < s[at] {
		return -1
	}
	return +1
}

func (s Slice[T]) Sort() {
	if !sort.IsSorted(s) {
		sort.Sort(s)
	}
}

func (s Slice[T]) FindLinear(needle T) (int, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == needle {
			return i, true
		}
	}
	return -1, false
}

func (s Slice[T]) FindBinary(needle T) (int, bool) {
	return sort.Find(len(s), func(i int) int { return s.Compare(needle, i) })
}
