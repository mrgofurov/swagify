package router

import (
	"encoding/json"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/swagify/core"
	"github.com/swagify/openapi"
	"github.com/swagify/ui"
)

// FiberAdapter integrates swagify with the Fiber web framework.
type FiberAdapter struct {
	app      *fiber.App
	registry *Registry
	group    fiber.Router

	// Configuration
	openAPIPath string
	docsPath    string

	// Docs authentication
	docsAuth *DocsAuthConfig
}

// FiberConfig holds configuration for the Fiber adapter.
type FiberConfig struct {
	// Info sets the API metadata (title, version, description, etc.)
	Info *core.Info

	// Servers sets the server URLs for the OpenAPI document.
	Servers []core.Server

	// OpenAPIPath sets the path to serve the OpenAPI JSON document.
	// Default: /openapi.json
	OpenAPIPath string

	// DocsPath sets the path to serve the docs UI.
	// Default: /docs
	DocsPath string

	// SecuritySchemes defines available security schemes.
	SecuritySchemes map[string]core.SecurityScheme

	// GlobalSecurity sets security requirements applied to all routes.
	GlobalSecurity []map[string][]string
}

// NewFiber creates a new FiberAdapter with the given Fiber app.
func NewFiber(app *fiber.App, configs ...FiberConfig) *FiberAdapter {
	adapter := &FiberAdapter{
		app:         app,
		registry:    NewRegistry(),
		group:       app,
		openAPIPath: "/openapi.json",
		docsPath:    "/docs",
	}

	if len(configs) > 0 {
		cfg := configs[0]
		if cfg.Info != nil {
			adapter.registry.SetInfo(cfg.Info)
		}
		for _, s := range cfg.Servers {
			adapter.registry.AddServer(s)
		}
		if cfg.OpenAPIPath != "" {
			adapter.openAPIPath = cfg.OpenAPIPath
		}
		if cfg.DocsPath != "" {
			adapter.docsPath = cfg.DocsPath
		}
		for name, scheme := range cfg.SecuritySchemes {
			adapter.registry.AddSecurityScheme(name, scheme)
		}
		if cfg.GlobalSecurity != nil {
			adapter.registry.SetGlobalSecurity(cfg.GlobalSecurity)
		}
	}

	return adapter
}

// Registry returns the underlying route registry.
func (f *FiberAdapter) Registry() *Registry {
	return f.registry
}

// SetInfo sets the API info metadata.
func (f *FiberAdapter) SetInfo(info *core.Info) {
	f.registry.SetInfo(info)
}

// AddServer adds a server to the API documentation.
func (f *FiberAdapter) AddServer(url, description string) {
	f.registry.AddServer(core.Server{URL: url, Description: description})
}

// AddTag adds a tag with a description.
func (f *FiberAdapter) AddTag(name, description string) {
	f.registry.AddTag(core.Tag{Name: name, Description: description})
}

// AddSecurityScheme adds a named security scheme.
func (f *FiberAdapter) AddSecurityScheme(name string, scheme core.SecurityScheme) {
	f.registry.AddSecurityScheme(name, scheme)
}

// AddBearerAuth adds a bearer token authentication security scheme.
func (f *FiberAdapter) AddBearerAuth() {
	f.registry.AddSecurityScheme("bearerAuth", core.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "Bearer token authentication",
	})
}

// AddAPIKeyAuth adds an API key authentication security scheme.
func (f *FiberAdapter) AddAPIKeyAuth(name, in string) {
	f.registry.AddSecurityScheme("apiKeyAuth", core.SecurityScheme{
		Type:        "apiKey",
		Name:        name,
		In:          in,
		Description: "API key authentication",
	})
}

// AddBasicAuth adds HTTP basic authentication security scheme.
func (f *FiberAdapter) AddBasicAuth() {
	f.registry.AddSecurityScheme("basicAuth", core.SecurityScheme{
		Type:        "http",
		Scheme:      "basic",
		Description: "Basic HTTP authentication",
	})
}

// GET registers a GET route with documentation.
func (f *FiberAdapter) GET(path string, handler fiber.Handler, opts ...RouteOption) {
	f.registerRoute("GET", path, handler, nil, nil, opts)
}

// POST registers a POST route with documentation.
func (f *FiberAdapter) POST(path string, handler fiber.Handler, opts ...RouteOption) {
	f.registerRoute("POST", path, handler, nil, nil, opts)
}

// PUT registers a PUT route with documentation.
func (f *FiberAdapter) PUT(path string, handler fiber.Handler, opts ...RouteOption) {
	f.registerRoute("PUT", path, handler, nil, nil, opts)
}

// PATCH registers a PATCH route with documentation.
func (f *FiberAdapter) PATCH(path string, handler fiber.Handler, opts ...RouteOption) {
	f.registerRoute("PATCH", path, handler, nil, nil, opts)
}

// DELETE registers a DELETE route with documentation.
func (f *FiberAdapter) DELETE(path string, handler fiber.Handler, opts ...RouteOption) {
	f.registerRoute("DELETE", path, handler, nil, nil, opts)
}

// TypedGET registers a typed GET route where request/response types are inferred.
func TypedGET[Req any, Res any](f *FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *fiber.Ctx) error {
		var req Req
		if reqType.Kind() == reflect.Struct {
			if err := c.QueryParser(&req); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "Invalid query parameters: " + err.Error()})
			}
		}

		res, err := handler(c, req)
		if err != nil {
			return err
		}
		return c.JSON(res)
	}

	f.registerRoute("GET", path, wrappedHandler, reqType, resType, opts)
}

// TypedPOST registers a typed POST route where request/response types are inferred.
func TypedPOST[Req any, Res any](f *FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *fiber.Ctx) error {
		var req Req
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
		}

		res, err := handler(c, req)
		if err != nil {
			return err
		}
		return c.Status(201).JSON(res)
	}

	f.registerRoute("POST", path, wrappedHandler, reqType, resType, opts)
}

// TypedPUT registers a typed PUT route where request/response types are inferred.
func TypedPUT[Req any, Res any](f *FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *fiber.Ctx) error {
		var req Req
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
		}

		res, err := handler(c, req)
		if err != nil {
			return err
		}
		return c.JSON(res)
	}

	f.registerRoute("PUT", path, wrappedHandler, reqType, resType, opts)
}

// TypedPATCH registers a typed PATCH route where request/response types are inferred.
func TypedPATCH[Req any, Res any](f *FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *fiber.Ctx) error {
		var req Req
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
		}

		res, err := handler(c, req)
		if err != nil {
			return err
		}
		return c.JSON(res)
	}

	f.registerRoute("PATCH", path, wrappedHandler, reqType, resType, opts)
}

// TypedDELETE registers a typed DELETE route where request/response types are inferred.
func TypedDELETE[Req any, Res any](f *FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *fiber.Ctx) error {
		var req Req
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
		}

		res, err := handler(c, req)
		if err != nil {
			return err
		}
		return c.Status(204).JSON(res)
	}

	f.registerRoute("DELETE", path, wrappedHandler, reqType, resType, opts)
}

// registerRoute creates a Route, applies options, registers it, and binds the handler.
func (f *FiberAdapter) registerRoute(method, path string, handler fiber.Handler, reqType, resType reflect.Type, opts []RouteOption) {
	route := &Route{
		Method:       method,
		Path:         path,
		RequestType:  reqType,
		ResponseType: resType,
	}

	applyOptions(route, opts)
	f.registry.Register(route)

	// Register with Fiber
	switch method {
	case "GET":
		f.group.Get(path, handler)
	case "POST":
		f.group.Post(path, handler)
	case "PUT":
		f.group.Put(path, handler)
	case "PATCH":
		f.group.Patch(path, handler)
	case "DELETE":
		f.group.Delete(path, handler)
	}
}

// BasicAuth protects the docs UI and OpenAPI JSON endpoints with HTTP Basic Authentication.
// Must be called before RegisterOpenAPI() and RegisterDocs().
//
// Simple usage:
//
//	api.BasicAuth("admin", "secret123")
//
// With config:
//
//	api.BasicAuth("admin", "secret123", swagify.DocsAuthConfig{Realm: "My API Docs"})
func (f *FiberAdapter) BasicAuth(username, password string, configs ...DocsAuthConfig) {
	cfg := DocsAuthConfig{}
	if len(configs) > 0 {
		cfg = configs[0]
	}
	cfg.Username = username
	cfg.Password = password
	f.docsAuth = &cfg
}

// RegisterOpenAPI registers the OpenAPI JSON endpoint.
func (f *FiberAdapter) RegisterOpenAPI(path ...string) {
	p := f.openAPIPath
	if len(path) > 0 {
		p = path[0]
		f.openAPIPath = p
	}

	handler := func(c *fiber.Ctx) error {
		gen := openapi.NewGenerator(f.registry)
		doc := gen.Generate()
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate OpenAPI document"})
		}
		c.Set("Content-Type", "application/json")
		return c.Send(data)
	}

	if f.docsAuth != nil {
		f.app.Get(p, fiberBasicAuth(*f.docsAuth), handler)
	} else {
		f.app.Get(p, handler)
	}
}

// RegisterDocs registers the docs UI endpoint.
func (f *FiberAdapter) RegisterDocs(path ...string) {
	p := f.docsPath
	if len(path) > 0 {
		p = path[0]
	}

	if f.docsAuth != nil {
		ui.RegisterFiberWithAuth(f.app, p, f.openAPIPath, fiberBasicAuth(*f.docsAuth))
	} else {
		ui.RegisterFiber(f.app, p, f.openAPIPath)
	}
}

// WithRequest sets the request body type for documentation on untyped handlers.
func WithRequest(model any) RouteOption {
	return func(r *Route) {
		if model != nil {
			r.RequestType = reflect.TypeOf(model)
		}
	}
}

// WithResponse sets the response body type for documentation on untyped handlers.
func WithResponse(model any) RouteOption {
	return func(r *Route) {
		if model != nil {
			r.ResponseType = reflect.TypeOf(model)
		}
	}
}

// Discover scans all existing routes registered on the Fiber app and
// automatically generates documentation entries for them. This allows
// using Swagify with existing projects without migrating route registration.
//
// Usage:
//
//	app := fiber.New()
//	app.Get("/users", listUsers)
//	app.Post("/users", createUser)
//
//	api := swagify.NewFiber(app)
//	api.Discover() // auto-documents all routes
//	api.RegisterOpenAPI()
//	api.RegisterDocs()
func (f *FiberAdapter) Discover(opts ...DiscoverOptions) {
	opt := DiscoverOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	// Set defaults
	if opt.AutoTags == nil {
		opt.AutoTags = boolPtr(true)
	}
	if opt.AutoSummary == nil {
		opt.AutoSummary = boolPtr(true)
	}

	// Use Fiber's GetRoutes to discover all registered routes.
	// The false parameter tells Fiber not to filter by method.
	routes := f.app.GetRoutes(true)

	// Deduplicate: Fiber may return duplicate routes for the same method+path
	seen := make(map[string]bool)

	for _, fr := range routes {
		method := fr.Method
		path := fr.Path

		// Skip unsupported methods
		if method == "HEAD" || method == "OPTIONS" || method == "TRACE" || method == "CONNECT" {
			continue
		}

		// Deduplicate
		key := method + " " + path
		if seen[key] {
			continue
		}
		seen[key] = true

		// Apply path filters
		if !shouldInclude(path, opt) {
			continue
		}

		// Check if this route is already registered (avoid duplicates from mixed usage)
		if existing := f.registry.FindRoute(key); existing != nil {
			continue
		}

		// Build the route
		route := &Route{
			Method: method,
			Path:   path,
		}

		// Auto-generate summary
		if *opt.AutoSummary {
			route.Summary = autoSummary(method, path)
		}

		// Auto-generate tags
		if *opt.AutoTags {
			tag := autoTag(path)
			if tag != "" {
				route.Tags = []string{tag}
			}
		}

		// Register (docs only — handler is already bound to the Fiber app)
		f.registry.Register(route)
	}
}

// Enrich adds metadata to a previously discovered (or registered) route.
// The key format is "METHOD /path", e.g., "GET /users/:id".
//
// Usage:
//
//	api.Discover()
//	api.Enrich("GET /users", swagify.Summary("List all users"), swagify.WithResponse(UsersResponse{}))
//	api.Enrich("POST /users", swagify.WithRequest(CreateUserReq{}), swagify.WithResponse(UserResponse{}))
func (f *FiberAdapter) Enrich(key string, opts ...RouteOption) {
	route := f.registry.FindRoute(key)
	if route == nil {
		return
	}

	applyOptions(route, opts)

	// Re-generate schemas for any new types added via enrichment
	if route.RequestType != nil {
		f.registry.SchemaGenerator().GenerateSchemaFromType(route.RequestType)
	}
	if route.ResponseType != nil {
		f.registry.SchemaGenerator().GenerateSchemaFromType(route.ResponseType)
	}
	if route.QueryType != nil {
		f.registry.SchemaGenerator().GenerateSchemaFromType(route.QueryType)
	}
	if route.PathType != nil {
		f.registry.SchemaGenerator().GenerateSchemaFromType(route.PathType)
	}
	if route.HeaderType != nil {
		f.registry.SchemaGenerator().GenerateSchemaFromType(route.HeaderType)
	}
	for _, resp := range route.Responses {
		if resp.Type != nil {
			f.registry.SchemaGenerator().GenerateSchemaFromType(resp.Type)
		}
	}
}
