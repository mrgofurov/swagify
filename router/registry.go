package router

import (
	"reflect"
	"strings"
	"sync"

	"github.com/swagify/core"
	"github.com/swagify/internal/utils"
	"github.com/swagify/openapi"
)

// Registry holds all registered routes and provides OpenAPI generation data.
type Registry struct {
	mu     sync.RWMutex
	routes []*Route
	schema *core.SchemaGenerator

	// Global metadata
	info     *core.Info
	servers  []core.Server
	tags     []core.Tag
	security []core.SecurityScheme

	// Global security requirements
	globalSecurity []map[string][]string
}

// NewRegistry creates a new route Registry with a fresh schema generator.
func NewRegistry() *Registry {
	return &Registry{
		routes: make([]*Route, 0),
		schema: core.NewSchemaGenerator(),
	}
}

// SetInfo sets the API info metadata.
func (r *Registry) SetInfo(info *core.Info) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.info = info
}

// AddServer adds a server to the API documentation.
func (r *Registry) AddServer(server core.Server) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.servers = append(r.servers, server)
}

// AddTag adds a tag to the API documentation.
func (r *Registry) AddTag(tag core.Tag) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags = append(r.tags, tag)
}

// AddSecurityScheme adds a security scheme.
func (r *Registry) AddSecurityScheme(name string, scheme core.SecurityScheme) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.security = append(r.security, scheme)
}

// SetGlobalSecurity sets global security requirements.
func (r *Registry) SetGlobalSecurity(security []map[string][]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.globalSecurity = security
}

// Register adds a route to the registry.
func (r *Registry) Register(route *Route) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Set OpenAPI path
	if route.OpenAPIPath == "" {
		route.OpenAPIPath = utils.NormalizePath(route.Path)
	}

	// Generate operation ID if not set
	if route.OperationID == "" {
		route.OperationID = utils.GenerateOperationID(route.Method, route.OpenAPIPath)
	}

	// Pre-generate schemas for request/response types
	if route.RequestType != nil {
		r.schema.GenerateSchemaFromType(route.RequestType)
	}
	if route.ResponseType != nil {
		r.schema.GenerateSchemaFromType(route.ResponseType)
	}
	if route.QueryType != nil {
		r.schema.GenerateSchemaFromType(route.QueryType)
	}
	if route.PathType != nil {
		r.schema.GenerateSchemaFromType(route.PathType)
	}
	if route.HeaderType != nil {
		r.schema.GenerateSchemaFromType(route.HeaderType)
	}

	// Generate schemas for custom response types
	for _, resp := range route.Responses {
		if resp.Type != nil {
			r.schema.GenerateSchemaFromType(resp.Type)
		}
	}

	r.routes = append(r.routes, route)
}

// FindRoute locates a registered route by "METHOD /path" (e.g., "GET /users/:id").
// Returns nil if not found. Used by Enrich() to add metadata to discovered routes.
func (r *Registry) FindRoute(key string) *Route {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := strings.SplitN(key, " ", 2)
	if len(parts) != 2 {
		return nil
	}
	method := strings.ToUpper(parts[0])
	path := strings.TrimSpace(parts[1])

	for _, route := range r.routes {
		if route.Method == method && route.Path == path {
			return route
		}
	}
	return nil
}

// Routes returns openapi-compatible Route objects.
// This implements the openapi.RegistryProvider interface.
func (r *Registry) Routes() []*openapi.Route {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*openapi.Route, len(r.routes))
	for i, route := range r.routes {
		result[i] = convertToOpenAPIRoute(route)
	}
	return result
}

// InternalRoutes returns the internal routes (for testing/debugging).
func (r *Registry) InternalRoutes() []*Route {
	r.mu.RLock()
	defer r.mu.RUnlock()
	routes := make([]*Route, len(r.routes))
	copy(routes, r.routes)
	return routes
}

// SchemaGenerator returns the schema generator used by this registry.
func (r *Registry) SchemaGenerator() *core.SchemaGenerator {
	return r.schema
}

// Info returns the API info.
func (r *Registry) Info() *core.Info {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.info != nil {
		return r.info
	}
	return &core.Info{
		Title:   "API Documentation",
		Version: "1.0.0",
	}
}

// Servers returns the configured servers.
func (r *Registry) Servers() []core.Server {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.servers
}

// Tags returns the configured tags.
func (r *Registry) Tags() []core.Tag {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.tags
}

// SecuritySchemes returns security schemes.
func (r *Registry) SecuritySchemes() []core.SecurityScheme {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.security
}

// GlobalSecurity returns the global security requirements.
func (r *Registry) GlobalSecurity() []map[string][]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.globalSecurity
}

// convertToOpenAPIRoute converts an internal Route to an openapi.Route.
func convertToOpenAPIRoute(r *Route) *openapi.Route {
	or := &openapi.Route{
		Method:              r.Method,
		Path:                r.Path,
		OpenAPIPath:         r.OpenAPIPath,
		RequestType:         r.RequestType,
		ResponseType:        r.ResponseType,
		Summary:             r.Summary,
		Description:         r.Description,
		OperationID:         r.OperationID,
		Tags:                r.Tags,
		Deprecated:          r.Deprecated,
		SuccessStatus:       r.SuccessStatus,
		Security:            r.Security,
		RequestContentType:  r.RequestContentType,
		ResponseContentType: r.ResponseContentType,
		QueryType:           r.QueryType,
		PathType:            r.PathType,
		HeaderType:          r.HeaderType,
	}

	// Convert parameters
	or.Parameters = r.Parameters

	// Convert responses
	if len(r.Responses) > 0 {
		or.Responses = make(map[int]openapi.RouteResponse, len(r.Responses))
		for code, resp := range r.Responses {
			or.Responses[code] = openapi.RouteResponse{
				StatusCode:  resp.StatusCode,
				Description: resp.Description,
				Type:        resp.Type,
				ContentType: resp.ContentType,
			}
		}
	}

	return or
}

// generateParamsFromType generates parameters from a struct type.
func generateParamsFromType(gen *core.SchemaGenerator, t reflect.Type, location string) []core.Parameter {
	if t == nil {
		return nil
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	var params []core.Parameter
	fields := core.StructFields(reflect.New(t).Elem().Interface())
	for _, f := range fields {
		schema := gen.GenerateSchemaFromType(f.Type)

		desc := ""
		if d := f.Tag.Get("description"); d != "" {
			desc = d
		}

		required := !f.OmitEmpty && f.Type.Kind() != reflect.Ptr
		if location == "path" {
			required = true
		}

		var example any
		if ex := f.Tag.Get("example"); ex != "" {
			example = ex
		}

		params = append(params, core.Parameter{
			Name:        f.JSONName,
			In:          location,
			Description: desc,
			Required:    required,
			Schema:      schema,
			Example:     example,
		})
	}

	return params
}
