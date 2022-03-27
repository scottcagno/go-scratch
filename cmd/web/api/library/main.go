package main

import (
	"log"
	"net/http"

	"github.com/scottcagno/go-scratch/pkg/web/api/rest"
	"github.com/scottcagno/go-scratch/pkg/web/api/rest/example"
)

var books example.Books

func init() {
	books := example.DefaultBookResource
	for i := 0; i < 15; i++ {

	}
}

func main() {
	api := rest.NewAPIServer("/api/v1/")
	api.RegisterResource("books", new(example.BooksResource))
	log.Println(http.ListenAndServe(":8080", api))
}
