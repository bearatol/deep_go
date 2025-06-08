package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	var result string
	obj := reflect.ValueOf(person)
	if obj.Kind() != reflect.Struct {
		return ""
	}

	objType := obj.Type()
	for i := range obj.NumField() {
		field := obj.Field(i)
		info := objType.Field(i)
		result += parseElement(&field, &info)
	}

	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	return result
}

func parseElement(field *reflect.Value, info *reflect.StructField) string {
	const reflectTag = "properties"
	elementTagList := info.Tag.Get(reflectTag)
	if len(elementTagList) == 0 {
		return ""
	}

	const emptyTag = "omitempty"
	tagParts := strings.Split(elementTagList, ",")
	if len(tagParts) > 2 {
		return ""
	}
	if tagParts[0] == emptyTag {
		return ""
	}

	value, is := isEmptyValue(field)
	if len(tagParts) == 2 && tagParts[1] == emptyTag && is {
		return ""
	}

	return fmt.Sprintf("%s=%s\n", tagParts[0], value)
}

func isEmptyValue(field *reflect.Value) (value string, isEmpty bool) {
	switch field.Kind() {
	case reflect.String:
		return field.String(), field.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", field.Int()), field.Int() == 0
	case reflect.Bool:
		return fmt.Sprintf("%t", field.Bool()), !field.Bool()
	default:
		return "", false
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
