package example

import (
	"sort"
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
	return sort.Search(len(b), func(i int) bool { return b[i].ID >= id })
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
