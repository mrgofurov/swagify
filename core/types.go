// Package core provides type reflection and JSON Schema generation
// for the swagify OpenAPI documentation package.
package core

// Schema represents a JSON Schema object compatible with OpenAPI 3.1.
type Schema struct {
	Type                 string            `json:"type,omitempty"`
	Format               string            `json:"format,omitempty"`
	Description          string            `json:"description,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Required             []string          `json:"required,omitempty"`
	Items                *Schema           `json:"items,omitempty"`
	AdditionalProperties *Schema           `json:"additionalProperties,omitempty"`
	Enum                 []any             `json:"enum,omitempty"`
	Example              any               `json:"example,omitempty"`
	Default              any               `json:"default,omitempty"`
	Nullable             bool              `json:"nullable,omitempty"`
	Ref                  string            `json:"$ref,omitempty"`
	Minimum              *float64          `json:"minimum,omitempty"`
	Maximum              *float64          `json:"maximum,omitempty"`
	MinLength            *int              `json:"minLength,omitempty"`
	MaxLength            *int              `json:"maxLength,omitempty"`
	Pattern              string            `json:"pattern,omitempty"`
	MinItems             *int              `json:"minItems,omitempty"`
	MaxItems             *int              `json:"maxItems,omitempty"`
	UniqueItems          bool              `json:"uniqueItems,omitempty"`
	Title                string            `json:"title,omitempty"`
	ReadOnly             bool              `json:"readOnly,omitempty"`
	WriteOnly            bool              `json:"writeOnly,omitempty"`
	Deprecated           bool              `json:"deprecated,omitempty"`
	AllOf                []*Schema         `json:"allOf,omitempty"`
	OneOf                []*Schema         `json:"oneOf,omitempty"`
	AnyOf                []*Schema         `json:"anyOf,omitempty"`
}

// Parameter represents an OpenAPI parameter (query, path, header, cookie).
type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"` // query, header, path, cookie
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
	Example     any     `json:"example,omitempty"`
	Deprecated  bool    `json:"deprecated,omitempty"`
}

// MediaType represents an OpenAPI media type object.
type MediaType struct {
	Schema  *Schema        `json:"schema,omitempty"`
	Example any            `json:"example,omitempty"`
	Examples map[string]any `json:"examples,omitempty"`
}

// RequestBody represents an OpenAPI request body.
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Required    bool                 `json:"required,omitempty"`
}

// Response represents an OpenAPI response object.
type Response struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
	Headers     map[string]*Schema   `json:"headers,omitempty"`
}

// SecurityScheme represents an OpenAPI security scheme.
type SecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	Name         string `json:"name,omitempty"`
	In           string `json:"in,omitempty"`
	Description  string `json:"description,omitempty"`
}

// Tag represents an OpenAPI tag.
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Info represents the OpenAPI info object.
type Info struct {
	Title          string  `json:"title"`
	Description    string  `json:"description,omitempty"`
	Version        string  `json:"version"`
	TermsOfService string  `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

// Contact represents the OpenAPI contact object.
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License represents the OpenAPI license object.
type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Server represents the OpenAPI server object.
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}
