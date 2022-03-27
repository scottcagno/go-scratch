package example

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/scottcagno/go-scratch/pkg/web/api/rest"
)

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
	// search book by id
	i := b.books.searchByID(id)
	// isolate book
	book := b.books[i]
	fn := func(w http.ResponseWriter, r *http.Request) {
		// return book
		rest.WriteAsJSON(w, book)
	}
	return http.HandlerFunc(fn)
}

func (b *BooksResource) Add(req *http.Request) http.Handler {
	// decode body into new book
	var book Book
	err := json.NewDecoder(req.Body).Decode(&book)
	// add book to set if no error
	if err == nil {
		b.books = append(b.books, book)
		// sort
		sort.Stable(b.books)
	}
	// now we can start handing...
	fn := func(w http.ResponseWriter, r *http.Request) {
		// if err exists
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// return books
		rest.WriteAsJSON(w, b.books)
	}
	return http.HandlerFunc(fn)
}

func (b *BooksResource) Set(req *http.Request, id string) http.Handler {
	// decode body into new book
	var book Book
	err := json.NewDecoder(req.Body).Decode(&book)
	// if no error
	if err == nil {
		// delete "old" book (for "update")
		delBookByID(&b.books, id)
		// and add new book (to complete "update")
		b.books = append(b.books, book)
		// sort
		sort.Stable(b.books)
	}
	// now we can start handling...
	fn := func(w http.ResponseWriter, r *http.Request) {
		// if err exists
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// return books
		rest.WriteAsJSON(w, b.books)
	}
	return http.HandlerFunc(fn)
}

func (b *BooksResource) Del(id string) http.Handler {
	// delete book
	delBookByID(&b.books, id)
	// now we can start handling...
	fn := func(w http.ResponseWriter, r *http.Request) {
		// return books
		rest.WriteAsJSON(w, b.books)
	}
	return http.HandlerFunc(fn)
}
