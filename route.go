package myapi

import "net/http"

//Route describes a single route in the API.
type Route struct {
	path        string
	description string
	method      string
	handler     http.HandlerFunc
	props       PropertyGroup
}

//NewRoute takes in all neccesay parts of a route: path method and handler function.
// other fields can be set via other functions.
func NewRoute(path, method string, Handler http.HandlerFunc) Route {
	return Route{
		path:    path,
		method:  method,
		handler: Handler,
	}
}

//SetDescription sets the description field on the Route.
func (r *Route) SetDescription(description string) {
	r.description = description
}

//AddPropertyGroup adds an entire propertygroup to the route properties
// note: this will replace any current properties on the route.
func (r *Route) AddPropertyGroup(pg PropertyGroup) {
	r.props = pg
}

//AddProperty adds all provided properties to the route.
// this will throw an error is a property name is duplicated.
func (r *Route) AddProperty(p ...Property) error {
	return r.props.AddProperties(p...)
}
