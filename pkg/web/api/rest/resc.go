package rest

import (
	"net/http"
	"path"
)

type Resource interface {
	GetAll() http.Handler
	Get(id string) http.Handler
	Add(req *http.Request) http.Handler
	Set(req *http.Request, id string) http.Handler
	Del(id string) http.Handler
}

// resource is an internal representation wrapping a user supplied ResourceHandler
type resource struct {
	name string
	path string
	Resource
}

// checkID returns a boolean reporting true if a resource id can be identified
func (re *resource) checkID(uri string) bool {
	if len(uri) > 0 && uri[len(uri)-1] == '/' {
		uri = uri[:len(uri)-1]
	}
	i := len(uri) - 1
	for i >= 0 && uri[i] != '/' {
		i--
	}
	return uri[i+1:] != "" && uri[i+1:] != re.name
}

func (re *resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	hasID := re.checkID(r.RequestURI)
	var id string
	if hasID {
		id = path.Base(r.RequestURI)
	}
	switch r.Method {
	case http.MethodGet:
		if hasID {
			h = re.Get(id)
			return
		}
		h = re.GetAll()
		return
	case http.MethodPost:
		h = re.Add(r)
		return
	case http.MethodPut:
		if hasID {
			h = re.Set(r, id)
			return
		}
	case http.MethodDelete:
		if hasID {
			h = re.Del(id)
			return
		}
	}
	h.ServeHTTP(w, r)
}
