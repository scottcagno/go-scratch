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

// resourceHandler is an internal representation wrapping a user supplied ResourceHandler
type resourceHandler struct {
	name string
	path string
	re   ResourceHandler
}

// checkID returns a boolean reporting true if a resource id can be identified
func (rh *resourceHandler) checkID(uri string) bool {
	if len(uri) > 0 && uri[len(uri)-1] == '/' {
		uri = uri[:len(uri)-1]
	}
	i := len(uri) - 1
	for i >= 0 && uri[i] != '/' {
		i--
	}
	return uri[i+1:] != "" && uri[i+1:] != rh.name
}

func (rh *resourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hasID := rh.checkID(r.RequestURI)
	switch r.Method {
	case http.MethodGet:
		if hasID {
			rh.re.GetOne(w, r)
			return
		}
		rh.re.GetAll(w, r)
		return
	case http.MethodPost:
		rh.re.AddOne(w, r)
		return
	case http.MethodPut:
		if hasID {
			rh.re.SetOne(w, r)
			return
		}
	case http.MethodDelete:
		if hasID {
			rh.re.DelOne(w, r)
			return
		}
	}
}
