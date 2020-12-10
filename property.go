package myapi

import (
	"fmt"
	"reflect"
)

//Property represents a single property in a request body.
type Property struct {
	Name     string
	propType Type
	rules    []Rule
}

//NewProperty creates a property with a blank rule set.
func NewProperty(name string, kind reflect.Kind) Property {
	return Property{
		Name:  name,
		rules: []Rule{},
	}
}

//AddRules will take the rules provided and add them to the Property,
// checking if they are valid first. If not, it will print a msg stating
// it has been ignored.
func (p *Property) AddRules(rules ...Rule) {
	for _, r := range rules {
		err := r.rulevalidation(*p)
		if err == nil {
			p.rules = append(p.rules, r)
		} else {
			fmt.Printf("Could not add rules to Property %v. error: %v \n", p.Name, err.Error())
		}
	}
}

//PropertyGroup wrapper used to be sure property names are unique when applied to a route.
type PropertyGroup struct {
	properties map[string]Property
}

//AddProperties attempts to add properties to PropertyGroup. It will throw an error if any Properties have
// conflicting names or aliases.
func (pg *PropertyGroup) AddProperties(props ...Property) error {
	for _, prop := range props {
		if _, present := pg.properties[prop.Name]; !present {
			pg.properties[prop.Name] = prop
		} else {
			return fmt.Errorf("duplicated Prop name: %v", prop.Name)
		}
	}
	return nil
}
