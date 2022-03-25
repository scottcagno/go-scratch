package library

// RID is a resource identifier
type RID string

// Author is an author resource
type Author struct {
	ID   RID    `json:"id,omitempty"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Genre is a genre resource
type Genre struct {
	ID   RID    `json:"id,omitempty"`
	Name string `json:"name"`
}

// Book is a book resource
type Book struct {
	ID     RID    `json:"id,omitempty"`
	Title  string `json:"title"`
	Author RID    `json:"author"`
	Genres []RID  `json:"genres"`
}
