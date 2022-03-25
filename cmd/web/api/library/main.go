package main

import (
	"fmt"
	"strings"
)

func main() {

	// mux := http.NewServeMux()
	// library.LibraryAPIRoutes(mux)
	// log.Panicln(http.ListenAndServe(":8080", mux))

	ss := []string{
		"/api/books",
		"/api/books/:id/",
		"/api/books/",
		"/api/books/:id",
		"/api/books/:id/author",
		"/api/books/:id/genres",
		"/api/author/:name/books",
	}
	for _, s := range ss {
		parsev2(s)
		// fmt.Printf("[%d] %#v\n", i, nss)
	}
}

func parse(s string) []string {
	ss := []string{s}
	var n int
	for {
		i := strings.IndexByte(s, '/')
		if i == -1 {
			return ss
		}
		n++
		s = s[i+1:]
		switch s[0] {
		case ':':
			ss = append(ss, s)
			i++
		}
		fmt.Println(s)
		// if i+1 != ':' {
		// 	s = s[i+1:]
		// 	fmt.Println("1>>>", s, i, string(s[i]))
		// 	continue
		// }
		// ss = append(ss, s)
		// fmt.Println("2>>>", s)
	}
	return ss
}

func parsev2(s string) {
	for {
		i := strings.IndexByte(s, '/')
		if i < 0 {
			break
		}
		s = s[i+1:]
		if len(s) > 0 && s[0] == ':' {
			fmt.Println(">>>>>>>>>>>>", s)
		}
		fmt.Printf(">>> s=%q, len(s)=%d, i=%d\n", s, len(s), i)
	}
}
