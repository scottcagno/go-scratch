package webkit

import (
	"context"
	"log"
	"net/http"
)

// Router plays nicely with the standard library's `net/http` package.
type Router interface {

	// A Handler responds to an HTTP request.
	http.Handler

	// Use appends one or more middlewares onto the Router stack.
	Use(middlewares ...func(http.Handler) http.Handler)

	// Handle and HandleFunc adds routes for `pattern` that matches
	// all HTTP methods.
	Handle(pattern string, h http.Handler)
	HandleFunc(pattern string, h http.HandlerFunc)

	// Method and MethodFunc adds routes for `pattern` that matches
	// the `method` HTTP method.
	Method(method, pattern string, h http.Handler)
	MethodFunc(method, pattern string, h http.HandlerFunc)

	// NotFound defines a handler to respond whenever a route could
	// not be found.
	NotFound(h http.HandlerFunc)

	// MethodNotAllowed defines a handler to respond whenever a method is
	// not allowed.
	MethodNotAllowed(h http.HandlerFunc)
}

func SetThenNext(k string, v any, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), k, v)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func GetThenNext(k string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := r.Context().Value(k)
			log.Printf("got value from context: key=%s, val=%v\n", k, v)
			next.ServeHTTP(w, r)
		},
	)
}
