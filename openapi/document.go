// Package openapi provides OpenAPI 3.1 document generation from
// the swagify route registry and schema generator.
package openapi

// Document represents a complete OpenAPI 3.1 document.
type Document struct {
	OpenAPI    string                `json:"openapi"`
	Info       InfoObject            `json:"info"`
	Servers    []ServerObject        `json:"servers,omitempty"`
	Paths      map[string]PathItem   `json:"paths"`
	Components *Components           `json:"components,omitempty"`
	Security   []map[string][]string `json:"security,omitempty"`
	Tags       []TagObject           `json:"tags,omitempty"`
}

// InfoObject represents the OpenAPI info section.
type InfoObject struct {
	Title          string         `json:"title"`
	Description    string         `json:"description,omitempty"`
	Version        string         `json:"version"`
	TermsOfService string         `json:"termsOfService,omitempty"`
	Contact        *ContactObject `json:"contact,omitempty"`
	License        *LicenseObject `json:"license,omitempty"`
}

// ContactObject represents the OpenAPI contact info.
type ContactObject struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// LicenseObject represents the OpenAPI license info.
type LicenseObject struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// ServerObject represents an OpenAPI server.
type ServerObject struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// TagObject represents an OpenAPI tag.
type TagObject struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// PathItem maps HTTP methods to their operations.
type PathItem map[string]*Operation

// Operation represents an OpenAPI operation (a single API endpoint).
type Operation struct {
	OperationID string                   `json:"operationId,omitempty"`
	Summary     string                   `json:"summary,omitempty"`
	Description string                   `json:"description,omitempty"`
	Tags        []string                 `json:"tags,omitempty"`
	Parameters  []ParameterObject        `json:"parameters,omitempty"`
	RequestBody *RequestBodyObject       `json:"requestBody,omitempty"`
	Responses   map[string]*ResponseObject `json:"responses"`
	Security    []map[string][]string    `json:"security,omitempty"`
	Deprecated  bool                     `json:"deprecated,omitempty"`
}

// ParameterObject represents an OpenAPI parameter.
type ParameterObject struct {
	Name        string       `json:"name"`
	In          string       `json:"in"`
	Description string       `json:"description,omitempty"`
	Required    bool         `json:"required,omitempty"`
	Schema      *SchemaObject `json:"schema,omitempty"`
	Example     any          `json:"example,omitempty"`
	Deprecated  bool         `json:"deprecated,omitempty"`
}

// RequestBodyObject represents an OpenAPI request body.
type RequestBodyObject struct {
	Description string                     `json:"description,omitempty"`
	Content     map[string]MediaTypeObject `json:"content"`
	Required    bool                       `json:"required,omitempty"`
}

// ResponseObject represents an OpenAPI response.
type ResponseObject struct {
	Description string                     `json:"description"`
	Content     map[string]MediaTypeObject `json:"content,omitempty"`
}

// MediaTypeObject represents an OpenAPI media type.
type MediaTypeObject struct {
	Schema *SchemaObject `json:"schema,omitempty"`
}

// SchemaObject is The OpenAPI Schema Object. It can be a full inline schema
// or a $ref pointer.
type SchemaObject struct {
	Type                 string                   `json:"type,omitempty"`
	Format               string                   `json:"format,omitempty"`
	Description          string                   `json:"description,omitempty"`
	Properties           map[string]*SchemaObject `json:"properties,omitempty"`
	Required             []string                 `json:"required,omitempty"`
	Items                *SchemaObject            `json:"items,omitempty"`
	AdditionalProperties *SchemaObject            `json:"additionalProperties,omitempty"`
	Enum                 []any                    `json:"enum,omitempty"`
	Example              any                      `json:"example,omitempty"`
	Default              any                      `json:"default,omitempty"`
	Nullable             bool                     `json:"nullable,omitempty"`
	Ref                  string                   `json:"$ref,omitempty"`
	Minimum              *float64                 `json:"minimum,omitempty"`
	Maximum              *float64                 `json:"maximum,omitempty"`
	MinLength            *int                     `json:"minLength,omitempty"`
	MaxLength            *int                     `json:"maxLength,omitempty"`
	Pattern              string                   `json:"pattern,omitempty"`
	MinItems             *int                     `json:"minItems,omitempty"`
	MaxItems             *int                     `json:"maxItems,omitempty"`
	UniqueItems          bool                     `json:"uniqueItems,omitempty"`
	Title                string                   `json:"title,omitempty"`
	ReadOnly             bool                     `json:"readOnly,omitempty"`
	WriteOnly            bool                     `json:"writeOnly,omitempty"`
	Deprecated           bool                     `json:"deprecated,omitempty"`
	AllOf                []*SchemaObject          `json:"allOf,omitempty"`
	OneOf                []*SchemaObject          `json:"oneOf,omitempty"`
	AnyOf                []*SchemaObject          `json:"anyOf,omitempty"`
}

// Components holds reusable OpenAPI components.
type Components struct {
	Schemas         map[string]*SchemaObject         `json:"schemas,omitempty"`
	SecuritySchemes map[string]*SecuritySchemeObject  `json:"securitySchemes,omitempty"`
}

// SecuritySchemeObject represents an OpenAPI security scheme.
type SecuritySchemeObject struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	Name         string `json:"name,omitempty"`
	In           string `json:"in,omitempty"`
	Description  string `json:"description,omitempty"`
}
