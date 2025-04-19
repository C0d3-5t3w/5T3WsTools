// Package reflect provides extensions to the standard library reflect package
package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// IsZero checks if a value is the zero value of its type
func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Bool:
		return !val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return val.Complex() == 0
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return val.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return val.IsNil()
	case reflect.Func:
		return val.IsNil()
	case reflect.Struct:
		return isZeroStruct(val)
	default:
		return false
	}
}

// isZeroStruct checks if all fields of a struct are zero values
func isZeroStruct(val reflect.Value) bool {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.CanInterface() && !IsZero(field.Interface()) {
			return false
		}
	}
	return true
}

// GetField safely gets a field from a struct by name, returning error if not found
func GetField(v interface{}, fieldName string) (interface{}, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct or pointer to struct")
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	if !field.CanInterface() {
		return nil, fmt.Errorf("cannot access field %s", fieldName)
	}

	return field.Interface(), nil
}

// SetField sets a field in a struct by name
func SetField(v interface{}, fieldName string, value interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer to struct")
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("not a pointer to struct")
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set field %s (unexported?)", fieldName)
	}

	valueVal := reflect.ValueOf(value)
	if !valueVal.Type().AssignableTo(field.Type()) {
		return fmt.Errorf("value type %s not assignable to field type %s",
			valueVal.Type().String(), field.Type().String())
	}

	field.Set(valueVal)
	return nil
}

// CallMethod calls a method on an object by name with given arguments
func CallMethod(v interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	val := reflect.ValueOf(v)
	method := val.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %s not found", methodName)
	}

	// Convert args to reflect.Values
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		reflectArgs[i] = reflect.ValueOf(arg)
	}

	// Call the method
	results := method.Call(reflectArgs)

	// Convert results back to interface{}
	ret := make([]interface{}, len(results))
	for i, result := range results {
		ret[i] = result.Interface()
	}

	return ret, nil
}

// StructToMap converts a struct to a map[string]interface{} with field names as keys
func StructToMap(v interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct or pointer to struct")
	}

	result := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.CanInterface() {
			// Get the field name, respecting json tag if present
			fieldName := typ.Field(i).Name
			tag := typ.Field(i).Tag.Get("json")
			if tag != "" && tag != "-" {
				parts := strings.Split(tag, ",")
				if parts[0] != "" {
					fieldName = parts[0]
				}
			}

			result[fieldName] = field.Interface()
		}
	}

	return result, nil
}

// MapToStruct fills a struct from a map[string]interface{}
func MapToStruct(m map[string]interface{}, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer to struct")
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("not a pointer to struct")
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldName := typ.Field(i).Name

		// Check JSON tag
		tag := typ.Field(i).Tag.Get("json")
		if tag != "" && tag != "-" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		// Find the value in the map
		mapVal, ok := m[fieldName]
		if !ok {
			continue
		}

		setErr := setFieldWithConversion(field, mapVal)
		if setErr != nil {
			return fmt.Errorf("error setting field %s: %w", fieldName, setErr)
		}
	}

	return nil
}

// setFieldWithConversion attempts to convert and set a field value
func setFieldWithConversion(field reflect.Value, value interface{}) error {
	valueVal := reflect.ValueOf(value)

	// Direct assignment if types are compatible
	if valueVal.Type().AssignableTo(field.Type()) {
		field.Set(valueVal)
		return nil
	}

	// Try type conversion for basic types
	switch field.Kind() {
	case reflect.String:
		field.SetString(fmt.Sprintf("%v", value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var intVal int64
		switch v := value.(type) {
		case int:
			intVal = int64(v)
		case int64:
			intVal = v
		case float64:
			intVal = int64(v)
		default:
			return fmt.Errorf("cannot convert %T to int", value)
		}
		field.SetInt(intVal)
	case reflect.Float32, reflect.Float64:
		var floatVal float64
		switch v := value.(type) {
		case int:
			floatVal = float64(v)
		case float64:
			floatVal = v
		default:
			return fmt.Errorf("cannot convert %T to float", value)
		}
		field.SetFloat(floatVal)
	default:
		return fmt.Errorf("unsupported conversion for field of type %s", field.Type().String())
	}

	return nil
}

// GetTag gets the value of a struct field tag
func GetTag(v interface{}, fieldName, tagName string) (string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("not a struct or pointer to struct")
	}

	typ := val.Type()
	field, ok := typ.FieldByName(fieldName)
	if !ok {
		return "", fmt.Errorf("field %s not found", fieldName)
	}

	return field.Tag.Get(tagName), nil
}

// HasMethod checks if a type has a specific method
func HasMethod(v interface{}, methodName string) bool {
	val := reflect.ValueOf(v)
	method := val.MethodByName(methodName)
	return method.IsValid()
}

// GetMethods returns all methods of a type
func GetMethods(v interface{}) []string {
	val := reflect.TypeOf(v)
	methods := make([]string, val.NumMethod())
	for i := 0; i < val.NumMethod(); i++ {
		methods[i] = val.Method(i).Name
	}
	return methods
}
