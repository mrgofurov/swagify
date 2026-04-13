// Package swagify provides a FastAPI-inspired HTTP framework for Go.
// Define typed handler functions and swagify automatically generates
// an OpenAPI specification and serves interactive documentation — no
// annotations required.
//
// Quick start:
//
//	api := swagify.New(swagify.Config{
//	    Title:   "My API",
//	    Version: "1.0.0",
//	})
//
//	api.GET("/users", listUsers)
//	api.POST("/users", createUser)
//
//	api.Run(":8080")
//	// → http://localhost:8080/docs      (Swagger UI)
//	// → http://localhost:8080/openapi.json
//
// Handler signatures supported:
//
//	func(http.ResponseWriter, *http.Request)   — plain net/http handler
//	func(*swagify.Ctx) error                   — access path/query/header via ctx
//	func(*swagify.Ctx) (Res, error)            — infers response schema
//	func(*swagify.Ctx, Req) (Res, error)       — infers request + response schemas
//
// For GET and DELETE routes, Req is parsed from URL query parameters.
// For POST, PUT, PATCH routes, Req is parsed from the JSON request body.
// Path parameters are always available via ctx.Param("name").
package swagify

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/mrgofurov/swagify/core"
	"github.com/mrgofurov/swagify/internal/utils"
	"github.com/mrgofurov/swagify/openapi"
	"github.com/mrgofurov/swagify/ui"
)

// Config holds the API configuration.
type Config struct {
	// Title is the API title shown in the docs. Default: "API Documentation".
	Title string

	// Description is the optional API description.
	Description string

	// Version is the API version string. Default: "1.0.0".
	Version string

	// DocsPath is the URL path for the interactive docs UI. Default: "/docs".
	DocsPath string

	// OpenAPIPath is the URL path for the OpenAPI JSON spec. Default: "/openapi.json".
	OpenAPIPath string

	// Servers lists the server URLs shown in the docs.
	Servers []core.Server
}

// API is the main entry point for swagify.
// Create one with New(), register routes, then call Run() or Handler().
type API struct {
	mux    *http.ServeMux
	mu     sync.RWMutex
	routes []*route
	schema *core.SchemaGenerator

	// Metadata
	info       core.Info
	servers    []core.Server
	tags       []core.Tag
	secSchemes []core.SecurityScheme
	globalSec  []map[string][]string

	// Paths
	docsPath    string
	openAPIPath string

	// Optional basic-auth protection for /docs and /openapi.json
	docsAuth *docsAuth
}

// docsAuth holds basic auth credentials for protecting the docs endpoints.
type docsAuth struct {
	username string
	password string
	realm    string
}

// route is the internal representation of a registered route.
type route struct {
	method      string
	path        string
	openAPIPath string

	reqType   reflect.Type
	resType   reflect.Type
	queryType reflect.Type

	summary     string
	description string
	operationID string
	tags        []string
	deprecated  bool

	successStatus int
	responses     map[int]routeResponse
	security      []map[string][]string
	reqCT         string
	resCT         string
}

// routeResponse documents a specific response status code.
type routeResponse struct {
	statusCode  int
	description string
	typ         reflect.Type
	contentType string
}

// New creates a new API with optional configuration.
func New(configs ...Config) *API {
	cfg := Config{
		Title:       "API Documentation",
		Version:     "1.0.0",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	}
	if len(configs) > 0 {
		c := configs[0]
		if c.Title != "" {
			cfg.Title = c.Title
		}
		if c.Description != "" {
			cfg.Description = c.Description
		}
		if c.Version != "" {
			cfg.Version = c.Version
		}
		if c.DocsPath != "" {
			cfg.DocsPath = c.DocsPath
		}
		if c.OpenAPIPath != "" {
			cfg.OpenAPIPath = c.OpenAPIPath
		}
		cfg.Servers = c.Servers
	}

	return &API{
		mux:    http.NewServeMux(),
		schema: core.NewSchemaGenerator(),
		info: core.Info{
			Title:       cfg.Title,
			Description: cfg.Description,
			Version:     cfg.Version,
		},
		servers:     cfg.Servers,
		docsPath:    cfg.DocsPath,
		openAPIPath: cfg.OpenAPIPath,
	}
}

// --- Route registration ---

// GET registers a GET route.
func (a *API) GET(path string, handler any, opts ...Option) {
	a.register("GET", path, handler, opts)
}

// POST registers a POST route.
func (a *API) POST(path string, handler any, opts ...Option) {
	a.register("POST", path, handler, opts)
}

// PUT registers a PUT route.
func (a *API) PUT(path string, handler any, opts ...Option) {
	a.register("PUT", path, handler, opts)
}

// PATCH registers a PATCH route.
func (a *API) PATCH(path string, handler any, opts ...Option) {
	a.register("PATCH", path, handler, opts)
}

// DELETE registers a DELETE route.
func (a *API) DELETE(path string, handler any, opts ...Option) {
	a.register("DELETE", path, handler, opts)
}

// register is the shared internal route registration.
func (a *API) register(method, path string, handler any, opts []Option) {
	result := scanHandler(method, handler)

	r := &route{
		method:    method,
		path:      path,
		reqType:   result.reqType,
		resType:   result.resType,
		queryType: result.queryType,
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.summary == "" {
		r.summary = autoSummary(method, path)
	}
	if len(r.tags) == 0 {
		if tag := autoTag(path); tag != "" {
			r.tags = []string{tag}
		}
	}

	r.openAPIPath = utils.NormalizePath(path)
	if r.operationID == "" {
		r.operationID = utils.GenerateOperationID(method, r.openAPIPath)
	}

	a.mu.Lock()
	if r.reqType != nil {
		a.schema.GenerateSchemaFromType(r.reqType)
	}
	if r.resType != nil {
		a.schema.GenerateSchemaFromType(r.resType)
	}
	if r.queryType != nil {
		a.schema.GenerateSchemaFromType(r.queryType)
	}
	for _, resp := range r.responses {
		if resp.typ != nil {
			a.schema.GenerateSchemaFromType(resp.typ)
		}
	}
	a.routes = append(a.routes, r)
	a.mu.Unlock()

	a.mux.HandleFunc(method+" "+path, result.handler)
}

// --- Server ---

// Run registers the docs endpoints and starts the HTTP server.
// It blocks until the server stops.
func (a *API) Run(addr string) error {
	a.registerSystemRoutes()
	return http.ListenAndServe(addr, a.mux)
}

// Handler registers the docs endpoints and returns the underlying http.Handler.
// Use this to embed swagify in an existing server.
func (a *API) Handler() http.Handler {
	a.registerSystemRoutes()
	return a.mux
}

// registerSystemRoutes mounts the OpenAPI JSON and docs UI endpoints.
func (a *API) registerSystemRoutes() {
	specHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gen := openapi.NewGenerator(a)
		doc := gen.Generate()
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			http.Error(w, "failed to generate spec", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data) //nolint:errcheck
	})

	if a.docsAuth != nil {
		mw := basicAuthMiddleware(a.docsAuth)
		a.mux.Handle("GET "+a.openAPIPath, mw(specHandler))
		ui.RegisterHTTPWithAuth(a.mux, a.docsPath, a.openAPIPath, mw)
	} else {
		a.mux.Handle("GET "+a.openAPIPath, specHandler)
		ui.RegisterHTTP(a.mux, a.docsPath, a.openAPIPath)
	}
}

// --- Metadata helpers ---

// AddTag adds a named tag (section) with an optional description to the docs.
func (a *API) AddTag(name, description string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.tags = append(a.tags, core.Tag{Name: name, Description: description})
}

// AddServer adds a server URL shown in the docs.
func (a *API) AddServer(url, description string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.servers = append(a.servers, core.Server{URL: url, Description: description})
}

// AddSecurityScheme registers a named security scheme.
func (a *API) AddSecurityScheme(name string, scheme core.SecurityScheme) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.secSchemes = append(a.secSchemes, scheme)
}

// AddBearerAuth registers a bearer token (JWT) security scheme.
func (a *API) AddBearerAuth() {
	a.AddSecurityScheme("bearerAuth", core.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "Bearer token authentication",
	})
}

// AddAPIKeyAuth registers an API key security scheme.
func (a *API) AddAPIKeyAuth(name, in string) {
	a.AddSecurityScheme("apiKeyAuth", core.SecurityScheme{
		Type:        "apiKey",
		Name:        name,
		In:          in,
		Description: "API key authentication",
	})
}

// BasicAuth protects the /docs and /openapi.json endpoints with HTTP basic
// authentication. Call this before Run() or Handler().
func (a *API) BasicAuth(username, password string, realm ...string) {
	r := "Swagify Docs"
	if len(realm) > 0 && realm[0] != "" {
		r = realm[0]
	}
	a.docsAuth = &docsAuth{username: username, password: password, realm: r}
}

// SetGlobalSecurity applies security requirements to every route.
func (a *API) SetGlobalSecurity(sec []map[string][]string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.globalSec = sec
}

// --- openapi.RegistryProvider implementation ---

func (a *API) Routes() []*openapi.Route {
	a.mu.RLock()
	defer a.mu.RUnlock()

	out := make([]*openapi.Route, len(a.routes))
	for i, r := range a.routes {
		or := &openapi.Route{
			Method:        r.method,
			Path:          r.path,
			OpenAPIPath:   r.openAPIPath,
			RequestType:   r.reqType,
			ResponseType:  r.resType,
			QueryType:     r.queryType,
			Summary:       r.summary,
			Description:   r.description,
			OperationID:   r.operationID,
			Tags:          r.tags,
			Deprecated:    r.deprecated,
			SuccessStatus: r.successStatus,
			Security:      r.security,
		}
		if r.reqCT != "" {
			or.RequestContentType = r.reqCT
		}
		if r.resCT != "" {
			or.ResponseContentType = r.resCT
		}
		if len(r.responses) > 0 {
			or.Responses = make(map[int]openapi.RouteResponse, len(r.responses))
			for code, resp := range r.responses {
				or.Responses[code] = openapi.RouteResponse{
					StatusCode:  resp.statusCode,
					Description: resp.description,
					Type:        resp.typ,
					ContentType: resp.contentType,
				}
			}
		}
		out[i] = or
	}
	return out
}

func (a *API) SchemaGenerator() *core.SchemaGenerator {
	return a.schema
}

func (a *API) Info() *core.Info {
	return &a.info
}

func (a *API) Servers() []core.Server {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.servers
}

func (a *API) Tags() []core.Tag {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.tags
}

func (a *API) SecuritySchemes() []core.SecurityScheme {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.secSchemes
}

func (a *API) GlobalSecurity() []map[string][]string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.globalSec
}

// --- Basic auth middleware ---

func basicAuthMiddleware(cfg *docsAuth) func(http.Handler) http.Handler {
	realm := cfg.realm
	expectedUser := sha256.Sum256([]byte(cfg.username))
	expectedPass := sha256.Sum256([]byte(cfg.password))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Basic ") {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(auth[6:])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			userHash := sha256.Sum256([]byte(parts[0]))
			passHash := sha256.Sum256([]byte(parts[1]))
			if subtle.ConstantTimeCompare(userHash[:], expectedUser[:]) != 1 ||
				subtle.ConstantTimeCompare(passHash[:], expectedPass[:]) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
