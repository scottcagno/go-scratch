package rest

import (
	"net/http"
)

type ResourceHandler interface {
	ReturnAll(w http.ResponseWriter, r *http.Request)
	ReturnOne(w http.ResponseWriter, r *http.Request)
	InsertOne(w http.ResponseWriter, r *http.Request)
	UpdateOne(w http.ResponseWriter, r *http.Request)
	DeleteOne(w http.ResponseWriter, r *http.Request)
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
