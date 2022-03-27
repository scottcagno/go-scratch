package library

import (
	"net/http"

	"github.com/scottcagno/go-scratch/pkg/web/api/rest"
)

// Author is an author resource
type Author struct {
	ID   rest.RID `json:"id,omitempty"`
	Name string   `json:"name"`
	Age  int      `json:"age"`
}

func (a Author) GetID() rest.RID {
	// TODO implement me
	panic("implement me")
}

func (a Author) SetID(id rest.RID) {
	// TODO implement me
	panic("implement me")
}

func (a Author) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// TODO implement me
	panic("implement me")
}
