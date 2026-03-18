package core

import (
	"reflect"
)

// ReflectType safely extracts the reflect.Type from a value,
// unwrapping pointers as needed.
func ReflectType(v any) reflect.Type {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// ReflectTypeName returns a clean name for a reflect.Type.
// Returns the struct name for named types, or a descriptive string for others.
func ReflectTypeName(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Name() != "" {
		return t.Name()
	}
	return t.String()
}

// IsStructType checks if a value represents a struct type (including pointers to structs).
func IsStructType(v any) bool {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

// StructFields returns the exported struct fields of a value with their json names.
func StructFields(v any) []FieldInfo {
	t := ReflectType(v)
	if t.Kind() != reflect.Struct {
		return nil
	}
	return extractFields(t)
}

// FieldInfo holds metadata about a struct field.
type FieldInfo struct {
	Name      string
	JSONName  string
	Type      reflect.Type
	Tag       reflect.StructTag
	OmitEmpty bool
	Ignored   bool
	Anonymous bool
}

func extractFields(t reflect.Type) []FieldInfo {
	var fields []FieldInfo
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		jsonTag := f.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		jsonName, omitempty := parseJSONTag(jsonTag)
		if jsonName == "" {
			jsonName = f.Name
		}

		fields = append(fields, FieldInfo{
			Name:      f.Name,
			JSONName:  jsonName,
			Type:      f.Type,
			Tag:       f.Tag,
			OmitEmpty: omitempty,
			Ignored:   jsonTag == "-",
			Anonymous: f.Anonymous,
		})
	}
	return fields
}
