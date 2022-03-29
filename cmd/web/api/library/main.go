package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/scottcagno/go-scratch/pkg/web/api/rest"
	"github.com/scottcagno/go-scratch/pkg/web/api/rest/example"
)

var books *example.BooksResource

func init() {
	books = example.DefaultBookResource
	for i := 0; i < 15; i++ {
		books.AppendBook(
			example.Book{
				ID:        fmt.Sprintf("%.3d", i),
				Title:     fmt.Sprintf("Book %.3d's Title", i),
				Author:    fmt.Sprintf("Book %.3d's Author", i),
				Published: time.Now().String(),
			},
		)
	}
}

func main() {
	api := rest.NewAPIServer("/api/v1/")
	api.RegisterResource("books", books)
	log.Println(http.ListenAndServe(":8080", api))
}
