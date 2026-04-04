package router

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/swagify/core"
	"github.com/swagify/openapi"
	"github.com/swagify/ui"
)

// HTTPAdapter integrates swagify with the standard net/http package.
type HTTPAdapter struct {
	mux      *http.ServeMux
	registry *Registry

	// Configuration
	openAPIPath string
	docsPath    string

	// Docs authentication
	docsAuth *DocsAuthConfig
}

// HTTPConfig holds configuration for the net/http adapter.
type HTTPConfig struct {
	// Info sets the API metadata.
	Info *core.Info

	// Servers sets the server URLs.
	Servers []core.Server

	// OpenAPIPath sets the path to serve the OpenAPI JSON document.
	OpenAPIPath string

	// DocsPath sets the path to serve the docs UI.
	DocsPath string
}

// NewHTTP creates a new HTTPAdapter with the given http.ServeMux.
func NewHTTP(mux *http.ServeMux, configs ...HTTPConfig) *HTTPAdapter {
	adapter := &HTTPAdapter{
		mux:         mux,
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
	}

	return adapter
}

// Registry returns the underlying route registry.
func (h *HTTPAdapter) Registry() *Registry {
	return h.registry
}

// SetInfo sets the API info metadata.
func (h *HTTPAdapter) SetInfo(info *core.Info) {
	h.registry.SetInfo(info)
}

// Handle registers a route with documentation.
func (h *HTTPAdapter) Handle(method, path string, handler http.HandlerFunc, opts ...RouteOption) {
	route := &Route{
		Method: method,
		Path:   path,
	}
	applyOptions(route, opts)
	h.registry.Register(route)

	// For Go 1.22+ pattern-based routing
	pattern := method + " " + path
	h.mux.HandleFunc(pattern, handler)
}

// GET registers a GET route.
func (h *HTTPAdapter) GET(path string, handler http.HandlerFunc, opts ...RouteOption) {
	h.Handle("GET", path, handler, opts...)
}

// POST registers a POST route.
func (h *HTTPAdapter) POST(path string, handler http.HandlerFunc, opts ...RouteOption) {
	h.Handle("POST", path, handler, opts...)
}

// PUT registers a PUT route.
func (h *HTTPAdapter) PUT(path string, handler http.HandlerFunc, opts ...RouteOption) {
	h.Handle("PUT", path, handler, opts...)
}

// PATCH registers a PATCH route.
func (h *HTTPAdapter) PATCH(path string, handler http.HandlerFunc, opts ...RouteOption) {
	h.Handle("PATCH", path, handler, opts...)
}

// DELETE registers a DELETE route.
func (h *HTTPAdapter) DELETE(path string, handler http.HandlerFunc, opts ...RouteOption) {
	h.Handle("DELETE", path, handler, opts...)
}

// HTTPTypedPOST registers a typed POST route for net/http.
func HTTPTypedPOST[Req any, Res any](h *HTTPAdapter, path string, handler func(r *http.Request, req Req) (Res, error), opts ...RouteOption) {
	var reqZero Req
	var resZero Res
	reqType := reflect.TypeOf(reqZero)
	resType := reflect.TypeOf(resZero)

	wrappedHandler := func(w http.ResponseWriter, r *http.Request) {
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body: " + err.Error()})
			return
		}

		res, err := handler(r, req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}

	route := &Route{
		Method:       "POST",
		Path:         path,
		RequestType:  reqType,
		ResponseType: resType,
	}
	applyOptions(route, opts)
	h.registry.Register(route)
	h.mux.HandleFunc("POST "+path, wrappedHandler)
}

// BasicAuth protects the docs UI and OpenAPI JSON endpoints with HTTP Basic Authentication.
// Must be called before RegisterOpenAPI() and RegisterDocs().
func (h *HTTPAdapter) BasicAuth(username, password string, configs ...DocsAuthConfig) {
	cfg := DocsAuthConfig{}
	if len(configs) > 0 {
		cfg = configs[0]
	}
	cfg.Username = username
	cfg.Password = password
	h.docsAuth = &cfg
}

// RegisterOpenAPI registers the OpenAPI JSON endpoint.
func (h *HTTPAdapter) RegisterOpenAPI(path ...string) {
	p := h.openAPIPath
	if len(path) > 0 {
		p = path[0]
		h.openAPIPath = p
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gen := openapi.NewGenerator(h.registry)
		doc := gen.Generate()
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	if h.docsAuth != nil {
		h.mux.Handle("GET "+p, httpBasicAuth(*h.docsAuth)(handler))
	} else {
		h.mux.Handle("GET "+p, handler)
	}
}

// RegisterDocs registers the docs UI endpoint.
func (h *HTTPAdapter) RegisterDocs(path ...string) {
	p := h.docsPath
	if len(path) > 0 {
		p = path[0]
	}

	if h.docsAuth != nil {
		ui.RegisterHTTPWithAuth(h.mux, p, h.openAPIPath, httpBasicAuth(*h.docsAuth))
	} else {
		ui.RegisterHTTP(h.mux, p, h.openAPIPath)
	}
}