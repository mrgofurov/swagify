package core

import (
	"reflect"
	"strings"
	"time"

	"github.com/mrgofurov/swagify/internal/cache"
)

// SchemaGenerator handles reflection-based JSON Schema generation
// with component deduplication and caching.
type SchemaGenerator struct {
	cache      *cache.SchemaCache
	components map[string]*Schema
}

// NewSchemaGenerator creates a new SchemaGenerator.
func NewSchemaGenerator() *SchemaGenerator {
	return &SchemaGenerator{
		cache:      cache.New(),
		components: make(map[string]*Schema),
	}
}

// GenerateSchema generates a JSON Schema from a Go value.
// Named struct types are registered as reusable components.
func (g *SchemaGenerator) GenerateSchema(v any) *Schema {
	if v == nil {
		return &Schema{Type: "object"}
	}
	t := reflect.TypeOf(v)
	return g.generateFromType(t)
}

// GenerateSchemaFromType generates a JSON Schema from a reflect.Type.
func (g *SchemaGenerator) GenerateSchemaFromType(t reflect.Type) *Schema {
	return g.generateFromType(t)
}

// Components returns all registered component schemas.
func (g *SchemaGenerator) Components() map[string]*Schema {
	result := make(map[string]*Schema, len(g.components))
	for k, v := range g.components {
		result[k] = v
	}
	return result
}

// SchemaRef returns a $ref schema pointing to a component, or the inline schema
// if the type is not a named struct.
func (g *SchemaGenerator) SchemaRef(v any) *Schema {
	if v == nil {
		return &Schema{Type: "object"}
	}
	t := reflect.TypeOf(v)
	return g.schemaRefFromType(t)
}

// SchemaRefFromType returns a $ref or inline schema from a reflect.Type.
func (g *SchemaGenerator) SchemaRefFromType(t reflect.Type) *Schema {
	return g.schemaRefFromType(t)
}

func (g *SchemaGenerator) schemaRefFromType(t reflect.Type) *Schema {
	// Unwrap pointer
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Named struct types get a component reference
	if t.Kind() == reflect.Struct && t.Name() != "" && !isTimeType(t) {
		name := g.typeName(t)
		// Ensure the component is generated
		if _, exists := g.components[name]; !exists {
			g.generateFromType(t)
		}
		return &Schema{Ref: "#/components/schemas/" + name}
	}

	return g.generateFromType(t)
}

func (g *SchemaGenerator) generateFromType(t reflect.Type) *Schema {
	// Unwrap pointer
	isPtr := false
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		isPtr = true
	}

	// Check cache for named types
	if t.Name() != "" {
		name := g.typeName(t)
		if cached, ok := g.cache.Get(name); ok {
			s := cached.(*Schema)
			if isPtr {
				copy := *s
				copy.Nullable = true
				return &copy
			}
			return s
		}
	}

	schema := g.resolveType(t)

	if isPtr && schema.Ref == "" {
		schema.Nullable = true
	}

	// Register named struct types as components
	if t.Kind() == reflect.Struct && t.Name() != "" && !isTimeType(t) {
		name := g.typeName(t)
		g.components[name] = schema
		g.cache.Set(name, schema)
	}

	return schema
}

func (g *SchemaGenerator) resolveType(t reflect.Type) *Schema {
	// Handle well-known types
	if isTimeType(t) {
		return &Schema{Type: "string", Format: "date-time"}
	}

	switch t.Kind() {
	case reflect.String:
		return &Schema{Type: "string"}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return &Schema{Type: "integer", Format: "int32"}

	case reflect.Int64:
		return &Schema{Type: "integer", Format: "int64"}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return &Schema{Type: "integer", Format: "int32"}

	case reflect.Uint64:
		return &Schema{Type: "integer", Format: "int64"}

	case reflect.Float32:
		return &Schema{Type: "number", Format: "float"}

	case reflect.Float64:
		return &Schema{Type: "number", Format: "double"}

	case reflect.Bool:
		return &Schema{Type: "boolean"}

	case reflect.Slice, reflect.Array:
		items := g.schemaRefFromType(t.Elem())
		return &Schema{Type: "array", Items: items}

	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			return &Schema{Type: "object"}
		}
		additional := g.schemaRefFromType(t.Elem())
		return &Schema{Type: "object", AdditionalProperties: additional}

	case reflect.Struct:
		return g.generateStructSchema(t)

	case reflect.Interface:
		return &Schema{Type: "object"}

	default:
		return &Schema{Type: "string"}
	}
}

func (g *SchemaGenerator) generateStructSchema(t reflect.Type) *Schema {
	schema := &Schema{
		Type:       "object",
		Properties: make(map[string]*Schema),
	}

	// Pre-register to handle recursive types
	if t.Name() != "" {
		name := g.typeName(t)
		g.components[name] = schema
		g.cache.Set(name, schema)
	}

	g.processStructFields(t, schema)
	return schema
}

func (g *SchemaGenerator) processStructFields(t reflect.Type, schema *Schema) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Handle embedded structs
		if field.Anonymous {
			embeddedType := field.Type
			for embeddedType.Kind() == reflect.Ptr {
				embeddedType = embeddedType.Elem()
			}
			if embeddedType.Kind() == reflect.Struct {
				g.processStructFields(embeddedType, schema)
				continue
			}
		}

		// Parse json tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		fieldName, omitempty := parseJSONTag(jsonTag)
		if fieldName == "" {
			fieldName = field.Name
		}

		// Generate field schema
		fieldSchema := g.generateFieldSchema(field)

		// Apply tag-based enrichments
		enrichSchemaFromTags(field, fieldSchema)

		schema.Properties[fieldName] = fieldSchema

		// Determine if field is required
		if isFieldRequired(field, omitempty) {
			schema.Required = append(schema.Required, fieldName)
		}
	}
}

func (g *SchemaGenerator) generateFieldSchema(field reflect.StructField) *Schema {
	ft := field.Type
	isPtr := false
	for ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
		isPtr = true
	}

	var schema *Schema

	// Named struct types use $ref
	if ft.Kind() == reflect.Struct && ft.Name() != "" && !isTimeType(ft) {
		schema = g.schemaRefFromType(ft)
	} else {
		schema = g.generateFromType(ft)
	}

	if isPtr && schema.Ref == "" {
		schema.Nullable = true
	}

	return schema
}

func (g *SchemaGenerator) typeName(t reflect.Type) string {
	name := t.Name()
	if name == "" {
		return "Anonymous"
	}
	return name
}

func isTimeType(t reflect.Type) bool {
	return t == reflect.TypeOf(time.Time{})
}

// parseJSONTag parses a json struct tag and returns the field name and whether omitempty is set.
func parseJSONTag(tag string) (string, bool) {
	if tag == "" {
		return "", false
	}
	parts := strings.Split(tag, ",")
	name := parts[0]
	omitempty := false
	for _, p := range parts[1:] {
		if p == "omitempty" {
			omitempty = true
		}
	}
	return name, omitempty
}
