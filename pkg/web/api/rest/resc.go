package rest

import (
	"net/http"
)

type Resource interface {

	// GetAll returns a http.Handler that locates and returns all
	// the implementing resource items.
	GetAll() http.Handler

	// GetOne takes an identifier and returns a http.Handler that
	// locates and returns the resource item with the matching
	// identifier.
	GetOne(id string) http.Handler

	// AddOne takes a serialized resource item (written to the request
	// body) and returns a http.Handler that adds the serialized item
	// to the resource set.
	AddOne(r *http.Request) http.Handler

	// SetOne takes an identifier along with a serialized resource
	// item (written to the request body) and returns a http.Handler
	// that locates and updates the resource item that has a matching
	// identifier.
	SetOne(r *http.Request, id string) http.Handler

	// DelOne takes an identifier and returns a http.Handler that
	// locates and deletes the resource item with the matching
	// identifier.
	DelOne(id string) http.Handler
}
