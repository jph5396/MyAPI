package myapi

import "reflect"

type (
	//Property represents a single property in a request body.
	Property struct {
		Name     string
		Alias    []string
		propKind reflect.Kind
		rules    []Rule
	}
)
