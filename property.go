package myapi

import (
	"errors"
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
func NewProperty(name string, typ Type) Property {
	return Property{
		Name:     name,
		propType: typ,
		rules:    []Rule{},
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

func (p *Property) validate(key string, value interface{}) error {
	valueType := reflect.TypeOf(value)
	//if value Type is a propertyGroup type,
	if p.propType == Group {
		return errors.New("PropGroup validation not yet implemented")
	}

	// make sure propType matches.
	if valueType != p.propType {
		return fmt.Errorf("%v: invalid type. got %v, want %v", key, valueType.String(), p.propType.String())
	}
	for _, rule := range p.rules {
		ok, msg := rule.validate(value)
		if !ok {
			return fmt.Errorf("%v: %v", key, msg)
		}
	}
	return nil
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

func (pg *PropertyGroup) validate(body map[string]interface{}) error {
	for key, val := range body {
		if property, ok := pg.properties[key]; ok {
			err := property.validate(key, val)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("property %v is not valid", key)
		}
	}

	return nil
}
