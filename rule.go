package myapi

import (
	"fmt"
	"reflect"
	"regexp"
)

//Rule Interface that defines the common interfaced that should be used when
// implementing a rule type all implemented rules can assume that the provided value
// is already an appropriate type, because the value will have been type checked before
// and invalid ones will not make it to the rule.
type Rule interface {
	//validate the function used to define if an input is valid or not.
	// if an input is not valid, the function should return false, and a string
	// representing the eror message that should be sent as part of the response.
	validate(interface{}) (bool, string)

	//rulevalidation should check to be sure the rule can be applied to given property.
	rulevalidation(Property) error
}

//RegexRule checks to see if the  propety value of a property matches the provided regex string.
type RegexRule struct {
	regexStr string
}

//NewRegexRule accepts a string, confirms it's a valid regex pattern, and returns a rule if the pattern is
//valid. If it is not, it will return a blank rule and an error noting that the regex is invalid.
func NewRegexRule(str string) (RegexRule, error) {

	_, err := regexp.Compile(str)
	if err != nil {
		return RegexRule{}, err
	}

	return RegexRule{
		regexStr: str,
	}, nil
}

func (r *RegexRule) validate(i interface{}) (bool, string) {

	value := i.(string)
	//Note: we ignore the error because it should have already been
	// confirmed that the expression string compiles when it was created.
	regex, _ := regexp.Compile(r.regexStr)
	if regex.MatchString(value) {
		return true, ""
	}
	msg := fmt.Sprintf("%v does not match regex pattern %v", value, r.regexStr)
	return false, msg
}

func (r *RegexRule) rulevalidation(p Property) error {
	if p.propKind != reflect.String {
		err := fmt.Errorf("regex rule cannot be used with property. got type %v, need string", p.propKind.String())
		return err
	}
	return nil
}

//EnumRule checks to see if the property value is within a set of valid values.
type EnumRule struct {
	enumKind   reflect.Kind
	enumvalues map[interface{}]struct{}
}

//NewEnumRule checks to see if all members provided are the same type,
// if so, it will return a valid EnumRule. if not, it will return
// a blank rule and a error. type is inferred by checking the first member of the array.
func NewEnumRule(members []interface{}) (EnumRule, error) {

	enumKind := reflect.TypeOf(members[0]).Kind()
	enumvalues := make(map[interface{}]struct{})
	var empty struct{}
	for _, val := range members {
		if reflect.TypeOf(val).Kind() != enumKind {
			err := fmt.Errorf("enum type mismatch")
			return EnumRule{}, err
		}
		enumvalues[val] = empty
	}
	return EnumRule{
		enumKind:   enumKind,
		enumvalues: enumvalues,
	}, nil
}

func (r *EnumRule) validate(i interface{}) (bool, string) {
	if _, ok := r.enumvalues[i]; ok {
		return true, ""
	}

	msg := fmt.Sprintf("%v not in enum list", i)
	return false, msg
}

func (r *EnumRule) rulevalidation(p Property) error {

	if r.enumKind != p.propKind {
		err := fmt.Errorf("enum type and prop type do not match. enum %v, prop type %v", r.enumKind.String(), p.propKind.String())
		return err
	}

	return nil
}

//CustomRule a rule type that can be used for user defined rules.
// see rule interface description for more information.
type CustomRule struct {
}
