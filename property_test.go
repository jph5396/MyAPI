package myapi

import (
	"testing"
)

func TestComplextProperty(t *testing.T) {
	// build properties.
	prop1 := NewProperty("Name", String)
	prop2 := NewProperty("ID", Int)
	prop3 := NewProperty("score", Float)

	propgroup := PropertyGroup{
		properties: make(map[string]Props),
	}
	propgroup.AddProperties(prop1, prop2, prop3)
	objProp := NewObjectProperty("User", false)
	objProp.UsePropertyGroup(propgroup)

	var test = map[string]interface{}{
		"User": map[string]interface{}{
			"Name":  "Jimbo",
			"ID":    43,
			"score": 23.45,
		},
		"UserFail": map[string]interface{}{
			"Name":  "Tomas",
			"ID":    43.4,
			"score": "test",
		},
	}

	err := objProp.validate("User", test["User"])
	if err != nil {
		t.Errorf("wanted nil, got %v", err.Error())
	}

	err = objProp.validate("UserFail", test["UserFail"])
	if err == nil {
		t.Errorf("wanted error got nil")
	}

}

func TestNestedProperties(t *testing.T) {
	prop1 := NewProperty("Name", String)
	prop2 := NewProperty("ID", Int)
	prop3 := NewProperty("score", Float)

	propgroup := PropertyGroup{
		properties: make(map[string]Props),
	}
	propgroup.AddProperties(prop1, prop2, prop3)
	objProp := NewObjectProperty("User", false)
	objProp.UsePropertyGroup(propgroup)

	computer := NewProperty("computer", String)
	status := NewProperty("status", String)
	supergroup := PropertyGroup{
		properties: make(map[string]Props),
	}
	supergroup.AddProperties(computer, status, objProp)
	superObjProp := NewObjectProperty("PCID", false)
	superObjProp.UsePropertyGroup(supergroup)
	var test = map[string]interface{}{
		"computer": "testPC",
		"status":   "working",
		"User": map[string]interface{}{
			"Name":  "Jimbo",
			"ID":    43,
			"score": 23.45,
		},
	}

	err := superObjProp.validate("nesteduser", test)
	if err != nil {
		t.Errorf("wanted nil got %v", err.Error())
	}

}

func TestArrayOfProperties(t *testing.T) {
	prop1 := NewProperty("Name", String)
	prop2 := NewProperty("ID", Int)
	prop3 := NewProperty("score", Float)

	propgroup := PropertyGroup{
		properties: make(map[string]Props),
	}
	propgroup.AddProperties(prop1, prop2, prop3)
	objProp := NewObjectProperty("Users", true)
	objProp.UsePropertyGroup(propgroup)

	var test = map[string]interface{}{
		"Users": []map[string]interface{}{
			{
				"Name":  "Jimbo",
				"ID":    43,
				"score": 23.45,
			},
			{
				"Name":  "Steven",
				"ID":    22,
				"score": 223.45,
			},
			{
				"Name":  "Paul",
				"ID":    11,
				"score": 2.45,
			},
		},
	}

	err := objProp.validate("list", test["Users"])
	if err != nil {
		t.Errorf("Wanted nil got %v", err.Error())
	}
}
