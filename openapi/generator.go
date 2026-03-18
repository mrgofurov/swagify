package openapi

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/swagify/core"
	"github.com/swagify/internal/utils"
)

// Route mirrors the router.Route type to avoid import cycles.
// The generator accepts this type directly.
type Route struct {
	Method       string
	Path         string
	OpenAPIPath  string
	RequestType  reflect.Type
	ResponseType reflect.Type

	Summary     string
	Description string
	OperationID string
	Tags        []string
	Deprecated  bool

	Parameters []core.Parameter
	Responses  map[int]RouteResponse

	SuccessStatus int
	Security      []map[string][]string

	RequestContentType  string
	ResponseContentType string

	QueryType  reflect.Type
	PathType   reflect.Type
	HeaderType reflect.Type
}

// RouteResponse represents a response definition.
type RouteResponse struct {
	StatusCode  int
	Description string
	Type        reflect.Type
	ContentType string
}

// DefaultSuccessStatus returns the appropriate default success status code.
func (r *Route) DefaultSuccessStatus() int {
	if r.SuccessStatus > 0 {
		return r.SuccessStatus
	}
	switch r.Method {
	case "POST":
		return 201
	case "DELETE":
		return 204
	default:
		return 200
	}
}

// DefaultRequestContentType returns the request content type.
func (r *Route) DefaultRequestContentType() string {
	if r.RequestContentType != "" {
		return r.RequestContentType
	}
	return "application/json"
}

// DefaultResponseContentType returns the response content type.
func (r *Route) DefaultResponseContentType() string {
	if r.ResponseContentType != "" {
		return r.ResponseContentType
	}
	return "application/json"
}

// RegistryProvider is the interface the generator needs from a route registry.
type RegistryProvider interface {
	Routes() []*Route
	SchemaGenerator() *core.SchemaGenerator
	Info() *core.Info
	Servers() []core.Server
	Tags() []core.Tag
	SecuritySchemes() []core.SecurityScheme
	GlobalSecurity() []map[string][]string
}

// Generator creates OpenAPI documents from a route registry.
type Generator struct {
	registry RegistryProvider
}

// NewGenerator creates a new OpenAPI Generator.
func NewGenerator(registry RegistryProvider) *Generator {
	return &Generator{registry: registry}
}

// Generate produces a complete OpenAPI 3.1 Document.
func (g *Generator) Generate() Document {
	info := g.registry.Info()
	gen := g.registry.SchemaGenerator()

	doc := Document{
		OpenAPI: "3.1.0",
		Info: InfoObject{
			Title:       info.Title,
			Description: info.Description,
			Version:     info.Version,
		},
		Paths: make(map[string]PathItem),
	}

	// Info extras
	if info.TermsOfService != "" {
		doc.Info.TermsOfService = info.TermsOfService
	}
	if info.Contact != nil {
		doc.Info.Contact = &ContactObject{
			Name:  info.Contact.Name,
			URL:   info.Contact.URL,
			Email: info.Contact.Email,
		}
	}
	if info.License != nil {
		doc.Info.License = &LicenseObject{
			Name: info.License.Name,
			URL:  info.License.URL,
		}
	}

	// Servers
	for _, s := range g.registry.Servers() {
		doc.Servers = append(doc.Servers, ServerObject{
			URL:         s.URL,
			Description: s.Description,
		})
	}

	// Tags
	for _, t := range g.registry.Tags() {
		doc.Tags = append(doc.Tags, TagObject{
			Name:        t.Name,
			Description: t.Description,
		})
	}

	// Build paths from routes
	for _, route := range g.registry.Routes() {
		openAPIPath := route.OpenAPIPath
		if openAPIPath == "" {
			openAPIPath = utils.NormalizePath(route.Path)
		}

		method := strings.ToLower(route.Method)
		operation := g.buildRouteOperation(route, gen)

		if _, exists := doc.Paths[openAPIPath]; !exists {
			doc.Paths[openAPIPath] = make(PathItem)
		}
		doc.Paths[openAPIPath][method] = operation
	}

	// Build components
	securitySchemes := buildSecuritySchemes(g.registry.SecuritySchemes())
	componentSchemas := gen.Components()
	doc.Components = buildComponents(componentSchemas, securitySchemes)

	// Collect unique tags from routes if not explicitly defined
	if len(doc.Tags) == 0 {
		tagSet := make(map[string]bool)
		for _, route := range g.registry.Routes() {
			for _, tag := range route.Tags {
				if !tagSet[tag] {
					tagSet[tag] = true
					doc.Tags = append(doc.Tags, TagObject{Name: tag})
				}
			}
		}
	}

	// Global security
	if globalSec := g.registry.GlobalSecurity(); len(globalSec) > 0 {
		doc.Security = globalSec
	}

	return doc
}

// buildRouteOperation converts a Route to an Operation.
func (g *Generator) buildRouteOperation(route *Route, gen *core.SchemaGenerator) *Operation {
	op := &Operation{
		OperationID: route.OperationID,
		Summary:     route.Summary,
		Description: route.Description,
		Tags:        route.Tags,
		Deprecated:  route.Deprecated,
		Responses:   make(map[string]*ResponseObject),
	}

	// Build path parameters from the path pattern
	pathParams := utils.ExtractPathParams(route.Path)
	for _, param := range pathParams {
		op.Parameters = append(op.Parameters, ParameterObject{
			Name:     param,
			In:       "path",
			Required: true,
			Schema:   &SchemaObject{Type: "string"},
		})
	}

	// Build query/path/header parameters from typed structs
	if route.QueryType != nil {
		params := g.buildParamsFromStructType(gen, route.QueryType, "query")
		op.Parameters = append(op.Parameters, params...)
	}
	if route.PathType != nil {
		params := g.buildParamsFromStructType(gen, route.PathType, "path")
		for _, tp := range params {
			replaced := false
			for i, existing := range op.Parameters {
				if existing.Name == tp.Name && existing.In == "path" {
					op.Parameters[i] = tp
					replaced = true
					break
				}
			}
			if !replaced {
				op.Parameters = append(op.Parameters, tp)
			}
		}
	}
	if route.HeaderType != nil {
		params := g.buildParamsFromStructType(gen, route.HeaderType, "header")
		op.Parameters = append(op.Parameters, params...)
	}

	// Build explicit parameters
	for _, p := range route.Parameters {
		op.Parameters = append(op.Parameters, ParameterObject{
			Name:        p.Name,
			In:          p.In,
			Description: p.Description,
			Required:    p.Required,
			Schema:      convertSchema(p.Schema),
			Example:     p.Example,
			Deprecated:  p.Deprecated,
		})
	}

	// Build request body
	if route.RequestType != nil && hasRequestBody(route.Method) {
		schema := gen.SchemaRefFromType(route.RequestType)
		contentType := route.DefaultRequestContentType()
		op.RequestBody = &RequestBodyObject{
			Required: true,
			Content: map[string]MediaTypeObject{
				contentType: {Schema: convertSchema(schema)},
			},
		}
	}

	// Build success response
	successStatus := route.DefaultSuccessStatus()
	successDesc := defaultStatusDescription(successStatus)
	statusStr := strconv.Itoa(successStatus)

	if route.ResponseType != nil {
		schema := gen.SchemaRefFromType(route.ResponseType)
		contentType := route.DefaultResponseContentType()
		op.Responses[statusStr] = &ResponseObject{
			Description: successDesc,
			Content: map[string]MediaTypeObject{
				contentType: {Schema: convertSchema(schema)},
			},
		}
	} else {
		op.Responses[statusStr] = &ResponseObject{
			Description: successDesc,
		}
	}

	// Build custom responses
	for status, resp := range route.Responses {
		sStr := strconv.Itoa(status)
		desc := resp.Description
		if desc == "" {
			desc = defaultStatusDescription(status)
		}
		if resp.Type != nil {
			schema := gen.SchemaRefFromType(resp.Type)
			ct := resp.ContentType
			if ct == "" {
				ct = "application/json"
			}
			op.Responses[sStr] = &ResponseObject{
				Description: desc,
				Content: map[string]MediaTypeObject{
					ct: {Schema: convertSchema(schema)},
				},
			}
		} else {
			op.Responses[sStr] = &ResponseObject{
				Description: desc,
			}
		}
	}

	// Security
	if len(route.Security) > 0 {
		op.Security = route.Security
	}

	return op
}

// buildParamsFromStructType generates OpenAPI parameters from a struct reflect.Type.
func (g *Generator) buildParamsFromStructType(gen *core.SchemaGenerator, t reflect.Type, location string) []ParameterObject {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	var params []ParameterObject

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		name := jsonTag
		if idx := strings.Index(name, ","); idx != -1 {
			name = name[:idx]
		}
		if name == "" {
			name = field.Name
		}

		fieldSchema := gen.GenerateSchemaFromType(field.Type)

		desc := field.Tag.Get("description")
		example := field.Tag.Get("example")

		required := location == "path"
		if !required {
			if validate := field.Tag.Get("validate"); strings.Contains(validate, "required") {
				required = true
			}
			if binding := field.Tag.Get("binding"); strings.Contains(binding, "required") {
				required = true
			}
		}

		p := ParameterObject{
			Name:        name,
			In:          location,
			Description: desc,
			Required:    required,
			Schema:      convertSchema(fieldSchema),
		}
		if example != "" {
			p.Example = example
		}

		params = append(params, p)
	}

	return params
}
