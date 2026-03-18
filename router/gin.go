package router

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/swagify/core"
	"github.com/swagify/openapi"
	"github.com/swagify/ui"
)

// GinAdapter integrates swagify with the Gin web framework.
type GinAdapter struct {
	engine   *gin.Engine
	group    *gin.RouterGroup
	registry *Registry

	// Configuration
	openAPIPath string
	docsPath    string
}

// GinConfig holds configuration for the Gin adapter.
type GinConfig struct {
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

// NewGin creates a new GinAdapter with the given Gin engine.
func NewGin(engine *gin.Engine, configs ...GinConfig) *GinAdapter {
	adapter := &GinAdapter{
		engine:      engine,
		group:       engine.Group(""),
		registry:    NewRegistry(),
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
func (g *GinAdapter) Registry() *Registry {
	return g.registry
}

// SetInfo sets the API info metadata.
func (g *GinAdapter) SetInfo(info *core.Info) {
	g.registry.SetInfo(info)
}

// AddServer adds a server to the API documentation.
func (g *GinAdapter) AddServer(url, description string) {
	g.registry.AddServer(core.Server{URL: url, Description: description})
}

// AddTag adds a tag with a description.
func (g *GinAdapter) AddTag(name, description string) {
	g.registry.AddTag(core.Tag{Name: name, Description: description})
}

// AddSecurityScheme adds a named security scheme.
func (g *GinAdapter) AddSecurityScheme(name string, scheme core.SecurityScheme) {
	g.registry.AddSecurityScheme(name, scheme)
}

// AddBearerAuth adds a bearer token authentication security scheme.
func (g *GinAdapter) AddBearerAuth() {
	g.registry.AddSecurityScheme("bearerAuth", core.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "Bearer token authentication",
	})
}

// AddAPIKeyAuth adds an API key authentication security scheme.
func (g *GinAdapter) AddAPIKeyAuth(name, in string) {
	g.registry.AddSecurityScheme("apiKeyAuth", core.SecurityScheme{
		Type:        "apiKey",
		Name:        name,
		In:          in,
		Description: "API key authentication",
	})
}

// AddBasicAuth adds HTTP basic authentication security scheme.
func (g *GinAdapter) AddBasicAuth() {
	g.registry.AddSecurityScheme("basicAuth", core.SecurityScheme{
		Type:        "http",
		Scheme:      "basic",
		Description: "Basic HTTP authentication",
	})
}

// GET registers a GET route with documentation.
func (g *GinAdapter) GET(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	g.registerRoute("GET", path, handler, nil, nil, opts)
}

// POST registers a POST route with documentation.
func (g *GinAdapter) POST(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	g.registerRoute("POST", path, handler, nil, nil, opts)
}

// PUT registers a PUT route with documentation.
func (g *GinAdapter) PUT(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	g.registerRoute("PUT", path, handler, nil, nil, opts)
}

// PATCH registers a PATCH route with documentation.
func (g *GinAdapter) PATCH(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	g.registerRoute("PATCH", path, handler, nil, nil, opts)
}

// DELETE registers a DELETE route with documentation.
func (g *GinAdapter) DELETE(path string, handler gin.HandlerFunc, opts ...RouteOption) {
	g.registerRoute("DELETE", path, handler, nil, nil, opts)
}

// GinTypedGET registers a typed GET route for Gin.
func GinTypedGET[Req any, Res any](g *GinAdapter, path string, handler func(*gin.Context, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *gin.Context) {
		var req Req
		if reqType.Kind() == reflect.Struct {
			if err := c.ShouldBindQuery(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
				return
			}
		}

		res, err := handler(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}

	g.registerRoute("GET", path, wrappedHandler, reqType, resType, opts)
}

// GinTypedPOST registers a typed POST route for Gin.
func GinTypedPOST[Req any, Res any](g *GinAdapter, path string, handler func(*gin.Context, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *gin.Context) {
		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		res, err := handler(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, res)
	}

	g.registerRoute("POST", path, wrappedHandler, reqType, resType, opts)
}

// GinTypedPUT registers a typed PUT route for Gin.
func GinTypedPUT[Req any, Res any](g *GinAdapter, path string, handler func(*gin.Context, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *gin.Context) {
		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		res, err := handler(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}

	g.registerRoute("PUT", path, wrappedHandler, reqType, resType, opts)
}

// GinTypedPATCH registers a typed PATCH route for Gin.
func GinTypedPATCH[Req any, Res any](g *GinAdapter, path string, handler func(*gin.Context, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *gin.Context) {
		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		res, err := handler(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}

	g.registerRoute("PATCH", path, wrappedHandler, reqType, resType, opts)
}

// GinTypedDELETE registers a typed DELETE route for Gin.
func GinTypedDELETE[Req any, Res any](g *GinAdapter, path string, handler func(*gin.Context, Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(c *gin.Context) {
		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		res, err := handler(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, res)
	}

	g.registerRoute("DELETE", path, wrappedHandler, reqType, resType, opts)
}

// registerRoute creates a Route, applies options, registers it, and binds the handler.
func (g *GinAdapter) registerRoute(method, path string, handler gin.HandlerFunc, reqType, resType reflect.Type, opts []RouteOption) {
	route := &Route{
		Method:       method,
		Path:         path,
		RequestType:  reqType,
		ResponseType: resType,
	}

	applyOptions(route, opts)
	g.registry.Register(route)

	// Register with Gin
	switch method {
	case "GET":
		g.group.GET(path, handler)
	case "POST":
		g.group.POST(path, handler)
	case "PUT":
		g.group.PUT(path, handler)
	case "PATCH":
		g.group.PATCH(path, handler)
	case "DELETE":
		g.group.DELETE(path, handler)
	}
}

// RegisterOpenAPI registers the OpenAPI JSON endpoint.
func (g *GinAdapter) RegisterOpenAPI(path ...string) {
	p := g.openAPIPath
	if len(path) > 0 {
		p = path[0]
		g.openAPIPath = p
	}

	g.engine.GET(p, func(c *gin.Context) {
		gen := openapi.NewGenerator(g.registry)
		doc := gen.Generate()
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OpenAPI document"})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	})
}

// RegisterDocs registers the docs UI endpoint.
func (g *GinAdapter) RegisterDocs(path ...string) {
	p := g.docsPath
	if len(path) > 0 {
		p = path[0]
	}
	ui.RegisterGin(g.engine, p, g.openAPIPath)
}