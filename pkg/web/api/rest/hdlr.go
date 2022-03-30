package rest

import (
	"net/http"
)

type ResourceHandler interface {

	// GetAll implements http.Handler and is responsible for locating
	// and returns all the implementing resource items.
	GetAll(w http.ResponseWriter, r *http.Request)

	// GetOne implements http.Handler and is responsible for locating
	// and returning the resource item with the matching identifier.
	// Note: the user is responsible for obtaining the identifier from
	// the request.
	GetOne(w http.ResponseWriter, r *http.Request)

	// AddOne implements http.Handler and is responsible for locating
	// the provided serialized resource item (written to the request
	// body) and adding the serialized item to the resource set.
	AddOne(w http.ResponseWriter, r *http.Request)

	// SetOne implements http.Handler and is responsible for locating
	// an identifier along with a serialized resource item (written to
	// the request body) and updating the resource item that has a
	// matching identifier.
	SetOne(w http.ResponseWriter, r *http.Request)

	// DelOne implements http.Handler and is responsible for locating
	// an identifier and deleting the resource item with the matching
	// identifier.
	DelOne(w http.ResponseWriter, r *http.Request)
}
