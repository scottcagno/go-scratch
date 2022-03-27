package rest

import (
	"net/http"
)

// RESTApiServer can be used to create API's
type RESTApiServer struct {
	base string
	h    *http.ServeMux
}

// NewAPIServer creates and returns a new server instance
func NewAPIServer(base string) *RESTApiServer {
	return &RESTApiServer{
		base: clean(base),
		h:    http.NewServeMux(),
	}
}

func (srv *RESTApiServer) RegisterResource(name string, re ResourceHandler) {
	rh := &resourceHandler{
		name: name,
		path: clean(srv.base + name),
		re:   re,
	}
	srv.h.Handle(rh.path, rh)
}

// ServeHTTP is the APIServer's default handler
func (srv *RESTApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// lookup resource handler
	h, pat := srv.h.Handler(r)
	// do something with pattern if we need to
	_ = pat
	// call handler
	h.ServeHTTP(w, r)
	return
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
			rh.re.ReturnOne(w, r)
			return
		}
		rh.re.ReturnAll(w, r)
		return
	case http.MethodPost:
		rh.re.InsertOne(w, r)
		return
	case http.MethodPut:
		if hasID {
			rh.re.UpdateOne(w, r)
			return
		}
	case http.MethodDelete:
		if hasID {
			rh.re.DeleteOne(w, r)
			return
		}
	}
}
