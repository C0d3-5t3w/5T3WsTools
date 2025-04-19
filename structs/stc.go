// Package structs provides utility functions for working with struct types
// that extend Go's standard library capabilities.
package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// ToMap converts a struct to a map[string]interface{}.
// Field tags can be used to customize the map keys.
func ToMap(s interface{}) (map[string]interface{}, error) {
	if s == nil {
		return nil, errors.New("input struct cannot be nil")
	}

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct or pointer to struct")
	}

	out := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Get the field name from the json tag if available
		name := field.Name
		if tag, ok := field.Tag.Lookup("json"); ok {
			name = parseTagName(tag)
			if name == "-" {
				continue
			}
		}

		out[name] = v.Field(i).Interface()
	}

	return out, nil
}

// FromMap converts a map[string]interface{} to a struct.
// Field tags can be used to customize the mapping.
func FromMap(m map[string]interface{}, s interface{}) error {
	if s == nil {
		return errors.New("output struct cannot be nil")
	}

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("output must be a pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("output must be a pointer to struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		name := field.Name
		if tag, ok := field.Tag.Lookup("json"); ok {
			tagName := parseTagName(tag)
			if tagName != "" && tagName != "-" {
				name = tagName
			}
		}

		if value, ok := m[name]; ok {
			fieldValue := v.Field(i)
			if fieldValue.CanSet() {
				val := reflect.ValueOf(value)
				// Only set if types are compatible
				if val.Type().AssignableTo(fieldValue.Type()) {
					fieldValue.Set(val)
				}
			}
		}
	}

	return nil
}

// Clone creates a deep copy of a struct
func Clone(src interface{}, dest interface{}) error {
	if src == nil || dest == nil {
		return errors.New("source and destination cannot be nil")
	}

	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("failed to marshal source: %w", err)
	}

	return json.Unmarshal(data, dest)
}

// Merge combines two structs. Values from src will overwrite values in dest if they exist.
// Both src and dest must be pointers to structs.
func Merge(src interface{}, dest interface{}) error {
	if src == nil || dest == nil {
		return errors.New("source and destination cannot be nil")
	}

	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if srcVal.Kind() != reflect.Ptr || destVal.Kind() != reflect.Ptr {
		return errors.New("both source and destination must be pointers")
	}

	srcVal = srcVal.Elem()
	destVal = destVal.Elem()

	if srcVal.Kind() != reflect.Struct || destVal.Kind() != reflect.Struct {
		return errors.New("both source and destination must be pointers to structs")
	}

	srcType := srcVal.Type()
	destType := destVal.Type()

	// First copy all fields with the same name and compatible types
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcType.Field(i)
		srcFieldVal := srcVal.Field(i)

		// Skip unexported fields
		if srcField.PkgPath != "" {
			continue
		}

		for j := 0; j < destVal.NumField(); j++ {
			destField := destType.Field(j)

			// Skip unexported fields
			if destField.PkgPath != "" {
				continue
			}

			if destField.Name == srcField.Name {
				destFieldVal := destVal.Field(j)
				if destFieldVal.CanSet() && srcFieldVal.Type().AssignableTo(destFieldVal.Type()) {
					destFieldVal.Set(srcFieldVal)
				}
				break
			}
		}
	}

	return nil
}

// HasField checks if a struct has a field with the given name
func HasField(s interface{}, fieldName string) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	_, found := v.Type().FieldByName(fieldName)
	return found
}

// GetField returns the value of a field from a struct by name
func GetField(s interface{}, fieldName string) (interface{}, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct or pointer to struct")
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s does not exist", fieldName)
	}

	return field.Interface(), nil
}

// SetField sets the value of a field in a struct by name
func SetField(s interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("input must be a pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("input must be a pointer to struct")
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s does not exist", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set (might be unexported)", fieldName)
	}

	val := reflect.ValueOf(value)
	if !val.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("value type %s is not assignable to field type %s",
			val.Type().String(), field.Type().String())
	}

	field.Set(val)
	return nil
}

// Helper function to parse tag names
func parseTagName(tag string) string {
	for i := 0; i < len(tag); i++ {
		if tag[i] == ',' {
			return tag[:i]
		}
	}
	return tag
}
