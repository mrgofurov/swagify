package swagify

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// parseQueryString parses URL query parameters into a struct pointer.
// Field names are resolved from the `json` struct tag, falling back to
// the lowercase field name.
func parseQueryString(r *http.Request, dst any) error {
	q := r.URL.Query()

	v := reflect.ValueOf(dst)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		name := jsonName(field)
		if name == "" {
			continue
		}

		vals, ok := q[name]
		if !ok || len(vals) == 0 {
			continue
		}

		fv := v.Field(i)
		ft := field.Type

		if ft.Kind() == reflect.Ptr {
			ptr := reflect.New(ft.Elem())
			if err := setScalar(ptr.Elem(), vals[0]); err == nil {
				fv.Set(ptr)
			}
		} else {
			setScalar(fv, vals[0]) //nolint:errcheck
		}
	}
	return nil
}

// jsonName extracts the JSON field name from a struct field tag.
func jsonName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return ""
	}
	if idx := strings.IndexByte(tag, ','); idx != -1 {
		tag = tag[:idx]
	}
	if tag == "" {
		return strings.ToLower(field.Name)
	}
	return tag
}

// setScalar sets a scalar reflect.Value from a string.
func setScalar(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(b)
	}
	return nil
}
