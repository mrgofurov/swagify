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
			scheme.Description = scheme.Description // preserve
			_ = name
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

// RegisterOpenAPI registers the OpenAPI JSON endpoint.
func (f *FiberAdapter) RegisterOpenAPI(path ...string) {
	p := f.openAPIPath
	if len(path) > 0 {
		p = path[0]
		f.openAPIPath = p
	}

	f.app.Get(p, func(c *fiber.Ctx) error {
		gen := openapi.NewGenerator(f.registry)
		doc := gen.Generate()
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate OpenAPI document"})
		}
		c.Set("Content-Type", "application/json")
		return c.Send(data)
	})
}

// RegisterDocs registers the docs UI endpoint.
func (f *FiberAdapter) RegisterDocs(path ...string) {
	p := f.docsPath
	if len(path) > 0 {
		p = path[0]
	}
	ui.RegisterFiber(f.app, p, f.openAPIPath)
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
