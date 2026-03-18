package core

// ExampleGenerator provides methods to generate example values for schemas.
type ExampleGenerator struct{}

// NewExampleGenerator creates a new ExampleGenerator.
func NewExampleGenerator() *ExampleGenerator {
	return &ExampleGenerator{}
}

// GenerateExample creates a basic example value for a schema type.
func (eg *ExampleGenerator) GenerateExample(schema *Schema) any {
	if schema.Example != nil {
		return schema.Example
	}

	if schema.Ref != "" {
		return nil
	}

	if len(schema.Enum) > 0 {
		return schema.Enum[0]
	}

	switch schema.Type {
	case "string":
		return eg.stringExample(schema)
	case "integer":
		return eg.integerExample(schema)
	case "number":
		return eg.numberExample(schema)
	case "boolean":
		return true
	case "array":
		if schema.Items != nil {
			item := eg.GenerateExample(schema.Items)
			if item != nil {
				return []any{item}
			}
		}
		return []any{}
	case "object":
		if schema.Properties != nil {
			obj := make(map[string]any)
			for name, prop := range schema.Properties {
				if ex := eg.GenerateExample(prop); ex != nil {
					obj[name] = ex
				}
			}
			return obj
		}
		return map[string]any{}
	}

	return nil
}

func (eg *ExampleGenerator) stringExample(schema *Schema) any {
	switch schema.Format {
	case "date-time":
		return "2024-01-01T00:00:00Z"
	case "date":
		return "2024-01-01"
	case "time":
		return "12:00:00"
	case "email":
		return "user@example.com"
	case "uri", "url":
		return "https://example.com"
	case "uuid":
		return "550e8400-e29b-41d4-a716-446655440000"
	case "ipv4":
		return "192.168.1.1"
	case "ipv6":
		return "::1"
	case "hostname":
		return "example.com"
	case "password":
		return "********"
	default:
		return "string"
	}
}

func (eg *ExampleGenerator) integerExample(schema *Schema) any {
	if schema.Minimum != nil {
		return int64(*schema.Minimum)
	}
	return int64(1)
}

func (eg *ExampleGenerator) numberExample(schema *Schema) any {
	if schema.Minimum != nil {
		return *schema.Minimum
	}
	return 1.0
}
