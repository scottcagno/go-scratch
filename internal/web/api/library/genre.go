package library

import (
	"github.com/scottcagno/go-scratch/pkg/web/api/rest"
)

// Genre is a genre resource
type Genre struct {
	ID   rest.RID `json:"id,omitempty"`
	Name string   `json:"name"`
}
