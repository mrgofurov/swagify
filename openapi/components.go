package openapi

import "github.com/mrgofurov/swagify/core"

// convertSchema converts an internal core.Schema to an OpenAPI SchemaObject.
func convertSchema(s *core.Schema) *SchemaObject {
	if s == nil {
		return nil
	}

	// If it's a $ref, return just the ref
	if s.Ref != "" {
		return &SchemaObject{Ref: s.Ref}
	}

	obj := &SchemaObject{
		Type:        s.Type,
		Format:      s.Format,
		Description: s.Description,
		Enum:        s.Enum,
		Example:     s.Example,
		Default:     s.Default,
		Nullable:    s.Nullable,
		Minimum:     s.Minimum,
		Maximum:     s.Maximum,
		MinLength:   s.MinLength,
		MaxLength:   s.MaxLength,
		Pattern:     s.Pattern,
		MinItems:    s.MinItems,
		MaxItems:    s.MaxItems,
		UniqueItems: s.UniqueItems,
		Title:       s.Title,
		ReadOnly:    s.ReadOnly,
		WriteOnly:   s.WriteOnly,
		Deprecated:  s.Deprecated,
		Required:    s.Required,
	}

	// Convert properties
	if len(s.Properties) > 0 {
		obj.Properties = make(map[string]*SchemaObject, len(s.Properties))
		for name, prop := range s.Properties {
			obj.Properties[name] = convertSchema(prop)
		}
	}

	// Convert items (array)
	if s.Items != nil {
		obj.Items = convertSchema(s.Items)
	}

	// Convert additionalProperties (map)
	if s.AdditionalProperties != nil {
		obj.AdditionalProperties = convertSchema(s.AdditionalProperties)
	}

	// Convert allOf
	if len(s.AllOf) > 0 {
		obj.AllOf = make([]*SchemaObject, len(s.AllOf))
		for i, a := range s.AllOf {
			obj.AllOf[i] = convertSchema(a)
		}
	}

	// Convert oneOf
	if len(s.OneOf) > 0 {
		obj.OneOf = make([]*SchemaObject, len(s.OneOf))
		for i, o := range s.OneOf {
			obj.OneOf[i] = convertSchema(o)
		}
	}

	// Convert anyOf
	if len(s.AnyOf) > 0 {
		obj.AnyOf = make([]*SchemaObject, len(s.AnyOf))
		for i, a := range s.AnyOf {
			obj.AnyOf[i] = convertSchema(a)
		}
	}

	return obj
}

// buildComponents creates the OpenAPI Components object from the registry.
func buildComponents(schemas map[string]*core.Schema, securitySchemes map[string]*SecuritySchemeObject) *Components {
	if len(schemas) == 0 && len(securitySchemes) == 0 {
		return nil
	}

	components := &Components{}

	if len(schemas) > 0 {
		components.Schemas = make(map[string]*SchemaObject, len(schemas))
		for name, schema := range schemas {
			components.Schemas[name] = convertSchema(schema)
		}
	}

	if len(securitySchemes) > 0 {
		components.SecuritySchemes = securitySchemes
	}

	return components
}
