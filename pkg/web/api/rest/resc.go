package rest

import (
	"log"
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

func LogRequest(r *http.Request, msg string) {
	log.Printf("method=%q, path=%q, msg=%q\n", r.Method, r.RequestURI, msg)
}

func (re *resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	hasID := re.checkID(r.RequestURI)
	var id string
	if hasID {
		id = path.Base(r.RequestURI)
	}
	if hasID {
		switch r.Method {
		case http.MethodGet:
			LogRequest(r, "get one")
			h = re.Get(id)
			goto serve
		case http.MethodPut:
			LogRequest(r, "update one")
			h = re.Set(r, id)
			goto serve
		case http.MethodDelete:
			LogRequest(r, "delete one")
			h = re.Del(id)
			goto serve
		default:
			LogRequest(r, "bad request with id")
			h = http.NotFoundHandler()
			goto serve
		}
	}
	switch r.Method {
	case http.MethodGet:
		LogRequest(r, "get all")
		h = re.GetAll()
		goto serve
	case http.MethodPost:
		LogRequest(r, "add one")
		h = re.Add(r)
		goto serve
	default:
		LogRequest(r, "bad request")
		h = http.NotFoundHandler()
		goto serve
	}
serve:
	h.ServeHTTP(w, r)
}
