package myapi

import "reflect"

//exposed funcs to make working with this package easier.

//PropsFromType receives a reflect.Type (generally of a struct) and
// returns a propertygroup based of the field name and types of the struct
// useful for creating propertygroups that dont need any specific rules applied to them.
func PropsFromType(t reflect.Type) PropertyGroup {

	return PropertyGroup{}
}
