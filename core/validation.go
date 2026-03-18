package core

import (
	"reflect"
	"strconv"
	"strings"
)

// enrichSchemaFromTags reads struct field tags and enriches the schema
// with validation constraints, formats, examples, and descriptions.
func enrichSchemaFromTags(field reflect.StructField, schema *Schema) {
	// Skip $ref schemas — they should not have inline properties
	if schema.Ref != "" {
		return
	}

	// description tag
	if desc := field.Tag.Get("description"); desc != "" {
		schema.Description = desc
	}

	// format tag
	if format := field.Tag.Get("format"); format != "" {
		schema.Format = format
	}

	// example tag
	if example := field.Tag.Get("example"); example != "" {
		schema.Example = parseExampleValue(example, schema.Type)
	}

	// deprecated tag
	if dep := field.Tag.Get("deprecated"); dep == "true" {
		schema.Deprecated = true
	}

	// readOnly tag
	if ro := field.Tag.Get("readOnly"); ro == "true" {
		schema.ReadOnly = true
	}

	// writeOnly tag
	if wo := field.Tag.Get("writeOnly"); wo == "true" {
		schema.WriteOnly = true
	}

	// Process validate tag (go-playground/validator style)
	if validate := field.Tag.Get("validate"); validate != "" {
		applyValidationTag(validate, schema)
	}

	// Process binding tag (gin style, similar format)
	if binding := field.Tag.Get("binding"); binding != "" {
		applyValidationTag(binding, schema)
	}

	// enum tag (comma-separated values)
	if enumTag := field.Tag.Get("enum"); enumTag != "" {
		values := strings.Split(enumTag, ",")
		enums := make([]any, len(values))
		for i, v := range values {
			enums[i] = strings.TrimSpace(v)
		}
		schema.Enum = enums
	}

	// default tag
	if def := field.Tag.Get("default"); def != "" {
		schema.Default = parseExampleValue(def, schema.Type)
	}
}

// applyValidationTag parses a validation tag string and applies
// constraints to the schema.
func applyValidationTag(tag string, schema *Schema) {
	rules := strings.Split(tag, ",")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		parts := strings.SplitN(rule, "=", 2)
		name := parts[0]
		value := ""
		if len(parts) == 2 {
			value = parts[1]
		}

		switch name {
		case "min":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				switch schema.Type {
				case "string":
					iv := int(v)
					schema.MinLength = &iv
				case "integer", "number":
					schema.Minimum = &v
				case "array":
					iv := int(v)
					schema.MinItems = &iv
				}
			}

		case "max":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				switch schema.Type {
				case "string":
					iv := int(v)
					schema.MaxLength = &iv
				case "integer", "number":
					schema.Maximum = &v
				case "array":
					iv := int(v)
					schema.MaxItems = &iv
				}
			}

		case "len":
			if v, err := strconv.Atoi(value); err == nil {
				schema.MinLength = &v
				schema.MaxLength = &v
			}

		case "email":
			schema.Format = "email"

		case "url", "uri":
			schema.Format = "uri"

		case "uuid", "uuid4":
			schema.Format = "uuid"

		case "ip":
			schema.Format = "ipv4"

		case "ipv4":
			schema.Format = "ipv4"

		case "ipv6":
			schema.Format = "ipv6"

		case "oneof":
			values := strings.Fields(value)
			enums := make([]any, len(values))
			for i, v := range values {
				enums[i] = v
			}
			schema.Enum = enums

		case "gt":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				exclusive := v
				schema.Minimum = &exclusive
			}

		case "gte":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				schema.Minimum = &v
			}

		case "lt":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				exclusive := v
				schema.Maximum = &exclusive
			}

		case "lte":
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				schema.Maximum = &v
			}
		}
	}
}

// isFieldRequired determines if a struct field should be marked as required
// based on json tags and validation tags.
func isFieldRequired(field reflect.StructField, omitempty bool) bool {
	// If omitempty is set, field is optional
	if omitempty {
		return false
	}

	// Pointer fields are optional by default
	if field.Type.Kind() == reflect.Ptr {
		return false
	}

	// Check validate tag for explicit "required"
	if validate := field.Tag.Get("validate"); validate != "" {
		rules := strings.Split(validate, ",")
		for _, rule := range rules {
			if strings.TrimSpace(rule) == "required" {
				return true
			}
		}
	}

	// Check binding tag for explicit "required"
	if binding := field.Tag.Get("binding"); binding != "" {
		rules := strings.Split(binding, ",")
		for _, rule := range rules {
			if strings.TrimSpace(rule) == "required" {
				return true
			}
		}
	}

	// Non-pointer, non-omitempty fields are required by default
	return true
}

// parseExampleValue attempts to parse an example string into the appropriate type.
func parseExampleValue(s string, schemaType string) any {
	switch schemaType {
	case "integer":
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return v
		}
	case "number":
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v
		}
	case "boolean":
		if v, err := strconv.ParseBool(s); err == nil {
			return v
		}
	}
	return s
}
