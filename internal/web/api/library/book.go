package library

import (
	"github.com/scottcagno/go-scratch/pkg/web/api/rest"
)

// Book is a book resource
type Book struct {
	ID     rest.RID   `json:"id,omitempty"`
	Title  string     `json:"title"`
	Author rest.RID   `json:"author"`
	Genres []rest.RID `json:"genres"`
}
