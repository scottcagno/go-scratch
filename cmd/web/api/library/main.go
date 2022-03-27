package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/scottcagno/go-scratch/pkg/web/api/rest"
)

type Book struct {
	ID    rest.RID `json:"id,omitempty"`
	Title string   `json:"title"`
}

func (b *Book) GetID() rest.RID {
	return b.ID
}

func (b *Book) SetID(id rest.RID) {
	b.ID = id
}

func (b *Book) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "book resource handler")
	return
}

func main() {
	api := rest.NewAPIServer("/api/v1/")
	api.HandleResource("book", new(Book))
	log.Panicln(http.ListenAndServe(":8080", api))
}

func handleAPITest() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		fmt.Fprintf(w, `{"id":"12345","title":"Moby Dick","cost":23.45,"author_id":"7643"}`)
	}
	return http.HandlerFunc(fn)
}

func handleFoo() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		fmt.Fprintf(w, "<h1>foo</h1>")
	}
	return http.HandlerFunc(fn)
}

func handleBar() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		fmt.Fprintf(w, "<h1>Bar</h1>")
	}
	return http.HandlerFunc(fn)
}

func matcherstuff() {
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
		// locate next '/'
		i := strings.IndexByte(s, '/')
		if i < 0 {
			break
		}
		//
		s = s[i+1:]
		if len(s) > 0 && s[0] == ':' {
			fmt.Println(">>>>>>>>>>>>", s)
		}
		fmt.Printf(">>> s=%q, len(s)=%d, i=%d\n", s, len(s), i)
	}
}
