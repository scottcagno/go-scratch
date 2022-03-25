package main

import (
	"path"
	"regexp"
	"strings"
	"testing"
)

var ss = []string{
	"/api/books",
	"/api/books/:id/",
	"/api/books/",
	"/api/books/:id",
	"/api/books/:id/author",
	"/api/books/:id/genres",
	"/api/author/:name/books",
}

func BenchmarkRawMatch(b *testing.B) {
	m := []string{
		"/api/books",
		"/api/books/:id/",
		"/api/books/",
		"/api/books/:id",
		"/api/books/:id/author",
		"/api/books/:id/genres",
		"/api/author/:name/books",
	}
	b.ReportAllocs()
	b.ResetTimer()
	var n int
	for i := 0; i < b.N; i++ {
		for _, s := range ss {
			for i := range m {
				if m[i] == s {
					n++
					continue
				}
			}
		}
		if n != len(ss) {
			b.Error("Fail")
		}
		n = 0
	}
}

func BenchmarkReplacerMatch(b *testing.B) {
	m := []string{
		"/api/books",
		"/api/books/123/",
		"/api/books/",
		"/api/books/123",
		"/api/books/123/author",
		"/api/books/123/genres",
		"/api/author/bobby/books",
	}
	b.ReportAllocs()
	b.ResetTimer()
	var n int
	for i := 0; i < b.N; i++ {
		for _, s := range ss {
			for i := range m {
				if strings.ReplaceAll(m[i], "123", ":id") == s {
					n++
					continue
				}
				if strings.ReplaceAll(m[i], "bobby", ":name") == s {
					n++
					continue
				}
			}
		}
		if n != len(ss) {
			b.Error("Fail")
		}
		n = 0
	}
}

func BenchmarkRegexMatch(b *testing.B) {
	m := regexp.MustCompile(`([A-z0-9_:\-/]+)`)
	b.ReportAllocs()
	b.ResetTimer()
	var n int
	for i := 0; i < b.N; i++ {
		for _, s := range ss {
			if m.MatchString(s) {
				n++
				continue
			}
		}
		if n != len(ss) {
			b.Error("Fail")
		}
		n = 0
	}
}

func BenchmarkGlobMatch(b *testing.B) {
	m := []string{
		"/api/books",
		"/api/books/*/*",
		"/api/books/",
		"/api/books/*",
	}
	b.ReportAllocs()
	b.ResetTimer()
	var n int
	for i := 0; i < b.N; i++ {
		for _, s := range ss {
			for i := range m {
				if ok, _ := path.Match(m[i], s); ok {
					n++
					continue
				}
			}
		}
		if n != len(ss) {
			b.Error("Fail")
		}
		n = 0
	}
}
