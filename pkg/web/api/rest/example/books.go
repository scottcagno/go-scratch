package example

import (
	"fmt"
	"net/http"
	"path"
	"sort"

	"github.com/scottcagno/go-scratch/pkg/web/api/rest"
)

type Book struct {
	ID        string `json:"id,omitempty"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Published string `json:"published"`
}

type Books []Book

func (b Books) Len() int {
	return len(b)
}

func (b Books) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Books) Less(i, j int) bool {
	return b[i].ID < b[j].ID
}

func (b Books) searchByID(id string) int {
	if !sort.IsSorted(b) {
		sort.Stable(b)
	}
	return sort.Search(len(b), func(i int) bool { return b[i].ID == id })
}

func (b Books) searchByTitle(title string) int {
	if !sort.IsSorted(b) {
		sort.Stable(b)
	}
	return sort.Search(len(b), func(i int) bool { return b[i].Title == title })
}

func (b Books) searchByAuthor(author string) int {
	if !sort.IsSorted(b) {
		sort.Stable(b)
	}
	return sort.Search(len(b), func(i int) bool { return b[i].Author == author })
}

func delBookByID(b *Books, id string) {
	// find book
	i := b.searchByID(id)
	// delete from slice (GC)
	if i < len(*b)-1 {
		copy((*b)[i:], (*b)[i+1:])
	}
	(*b)[len(*b)-1] = Book{} // or the zero value of T
	*b = (*b)[:len(*b)-1]
}

type BooksHandler struct {
	books []Book
}

func (b *BooksHandler) ReturnAll(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%q, return all (%q)\n", b, r.RequestURI)
}

func (b *BooksHandler) ReturnOne(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.RequestURI)
	fmt.Fprintf(w, "%q, return one (%q), id=%v\n", b, r.RequestURI, id)
}

func (b *BooksHandler) InsertOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%q, insert one\n", b)
}

func (b *BooksHandler) UpdateOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%q, update one\n", b)
}

func (b *BooksHandler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%q, delete one\n", b)
}

func (b *BooksHandler) String() string {
	return "BooksHandler"
}

type BooksResource struct {
	books Books
}

func (b *BooksResource) GetAll() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// return all books
		rest.WriteAsJSON(w, b.books)
	}
	return http.HandlerFunc(fn)
}

func (b *BooksResource) Get(id string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// search book by id
		i := b.books.searchByID(id)
		// return book
		rest.WriteAsJSON(w, b.books[i])
	}
	return http.HandlerFunc(fn)
}

func (b *BooksResource) Add() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// implement
		rest.WriteAsJSON(w, struct{ Msg string }{"Implement me..."})
	}
	return http.HandlerFunc(fn)
}

func (b *BooksResource) Set(id string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// implement
		rest.WriteAsJSON(w, struct{ Msg string }{"Implement me..."})
	}
	return http.HandlerFunc(fn)
}

func (b *BooksResource) Del(id string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// delete book
		delBookByID(&b.books, id)
		// return books
		rest.WriteAsJSON(w, b.books)
	}
	return http.HandlerFunc(fn)
}
