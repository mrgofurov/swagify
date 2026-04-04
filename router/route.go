// Package router provides route registration, metadata collection,
// and framework adapter interfaces for the swagify documentation package.
package router

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/mrgofurov/swagify/core"
)

// Route represents a registered API route with all its metadata.
type Route struct {
	// HTTP method (GET, POST, PUT, PATCH, DELETE)
	Method string

	// Original path as registered with the framework (e.g., /users/:id)
	Path string

	// Normalized OpenAPI path (e.g., /users/{id})
	OpenAPIPath string

	// Request body type (nil if no body)
	RequestType reflect.Type

	// Response body type (nil if no response body)
	ResponseType reflect.Type

	// Operation metadata
	Summary     string
	Description string
	OperationID string
	Tags        []string
	Deprecated  bool

	// Explicit parameters
	Parameters []core.Parameter

	// Custom responses (status code => response)
	Responses map[int]RouteResponse

	// Default success status code
	SuccessStatus int

	// Security requirements for this route
	Security []map[string][]string

	// Content type overrides
	RequestContentType  string
	ResponseContentType string

	// Query/Path/Header parameter types
	QueryType  reflect.Type
	PathType   reflect.Type
	HeaderType reflect.Type

	// Middlewares
	FiberMiddlewares []fiber.Handler
	GinMiddlewares   []gin.HandlerFunc
}

// RouteResponse represents a response definition for a specific status code.
type RouteResponse struct {
	StatusCode  int
	Description string
	Type        reflect.Type
	ContentType string
}

// DefaultSuccessStatus returns the appropriate default success status code
// based on the HTTP method if no explicit status is configured.
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

// DefaultRequestContentType returns the configured or default request content type.
func (r *Route) DefaultRequestContentType() string {
	if r.RequestContentType != "" {
		return r.RequestContentType
	}
	return "application/json"
}

// DefaultResponseContentType returns the configured or default response content type.
func (r *Route) DefaultResponseContentType() string {
	if r.ResponseContentType != "" {
		return r.ResponseContentType
	}
	return "application/json"
}
