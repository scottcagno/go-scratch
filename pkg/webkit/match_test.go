package webkit

import (
	"fmt"
	"path"
	"regexp"
	"testing"
)

type MatchTest struct {
	pattern, s string
	match      bool
	err        error
}

var matchTests = []MatchTest{
	{"abc", "abc", true, nil},
	{"*", "abc", true, nil},
	{"*c", "abc", true, nil},
	{"a*", "a", true, nil},
	{"a*", "abc", true, nil},
	{"a*", "ab/c", false, nil},
	{"a*/b", "abc/b", true, nil},
	{"a*/b", "a/c/b", false, nil},
	{"ab[c]", "abc", true, nil},
	{"ab[b-d]", "abc", true, nil},
	{"ab[e-g]", "abc", false, nil},
	{"ab[^e-g]", "abc", true, nil},
	{"a\\*b", "a*b", true, nil},
	{"a?b", "a☺b", true, nil},
	{"a[^a]b", "a☺b", true, nil},
	{"*x", "xxx", true, nil},
}

func TestMatch(t *testing.T) {
	for _, tt := range matchTests {
		ok, err := path.Match(tt.pattern, tt.s)
		if ok != tt.match || err != tt.err {
			t.Errorf("Match(%#q, %#q) = %v, %v want %v, %v", tt.pattern, tt.s, ok, err, tt.match, tt.err)
		}
	}
}

func BenchmarkMatch(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, tt := range matchTests {
			ok, err := path.Match(tt.pattern, tt.s)
			if ok != tt.match || err != tt.err {
				b.Errorf("Match(%#q, %#q) = %v, %v want %v, %v", tt.pattern, tt.s, ok, err, tt.match, tt.err)
			}
		}
	}
}

type FindTest struct {
	pat     string
	text    string
	matches [][]int
}

func (t FindTest) String() string {
	return fmt.Sprintf("pat: %#q text: %#q", t.pat, t.text)
}

var findTests = []FindTest{
	{`^abcdefg`, "abcdefg", build(1, 0, 7)},
	{`a+`, "baaab", build(1, 1, 4)},
	{"abcd..", "abcdef", build(1, 0, 6)},
	{`a`, "a", build(1, 0, 1)},
	{`b`, "abc", build(1, 1, 2)},
	{`.*`, "abcdef", build(1, 0, 6)},
	{`^`, "abcde", build(1, 0, 0)},
	{`$`, "abcde", build(1, 5, 5)},
	{`^abcd$`, "abcd", build(1, 0, 4)},
	{`[a-z]+`, "abcd", build(1, 0, 4)},
	{`[^a-z]+`, "ab1234cd", build(1, 2, 6)},
	{`a*(|(b))c*`, "aacc", build(1, 0, 4, 2, 2, -1, -1)},
	{`(.*).*`, "ab", build(1, 0, 2, 0, 2)},
	{`.`, "abc", build(3, 0, 1, 1, 2, 2, 3)},
	{`.(.)`, "abcd", build(2, 0, 2, 1, 2, 2, 4, 3, 4)},
	{`a(b*)`, "abbaab", build(3, 0, 3, 1, 3, 3, 4, 4, 4, 4, 6, 5, 6)},
}

// build is a helper to construct a [][]int by extracting n sequences from x.
// This represents n matches with len(x)/n submatches each.
func build(n int, x ...int) [][]int {
	ret := make([][]int, n)
	runLength := len(x) / n
	j := 0
	for i := range ret {
		ret[i] = make([]int, runLength)
		copy(ret[i], x[j:])
		j += runLength
		if j > len(x) {
			panic("invalid build entry")
		}
	}
	return ret
}

func TestFindString(t *testing.T) {
	for _, test := range findTests {
		result := regexp.MustCompile(test.pat).FindString(test.text)
		switch {
		case len(test.matches) == 0 && len(result) == 0:
			// ok
		case test.matches == nil && result != "":
			t.Errorf("expected no match; got one: %s", test)
		case test.matches != nil && result == "":
			// Tricky because an empty result has two meanings: no match or empty match.
			if test.matches[0][0] != test.matches[0][1] {
				t.Errorf("expected match; got none: %s", test)
			}
		case test.matches != nil && result != "":
			expect := test.text[test.matches[0][0]:test.matches[0][1]]
			if expect != result {
				t.Errorf("expected %q got %q: %s", expect, result, test)
			}
		}
	}
}

func Benchmark_Match_Regex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range findTests {
			result := regexp.MustCompile(test.pat).FindString(test.text)
			switch {
			case len(test.matches) == 0 && len(result) == 0:
				// ok
			case test.matches == nil && result != "":
				b.Errorf("expected no match; got one: %s", test)
			case test.matches != nil && result == "":
				// Tricky because an empty result has two meanings: no match or empty match.
				if test.matches[0][0] != test.matches[0][1] {
					b.Errorf("expected match; got none: %s", test)
				}
			case test.matches != nil && result != "":
				expect := test.text[test.matches[0][0]:test.matches[0][1]]
				if expect != result {
					b.Errorf("expected %q got %q: %s", expect, result, test)
				}
			}
		}
	}
}
