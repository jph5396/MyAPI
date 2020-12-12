package myapi

import "net/http"

//Middleware represents a middleware function that uses http.Handler
type Middleware struct {
	name        string
	description string
	handler     func(http.Handler) http.Handler
}

//NewMiddleware returns a new middleware with the name and handler provided.
func NewMiddleware(name string, handler func(http.Handler) http.Handler) Middleware {
	return Middleware{name: name, handler: handler}
}

//SetDescription set the description of the middleware
func (m *Middleware) SetDescription(desc string) {
	m.description = desc
}
