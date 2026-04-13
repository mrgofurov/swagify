package swagify

import (
	"context"
	"net/http"
)

// Ctx is the request context passed to typed handlers.
// It provides access to path parameters, headers, and the underlying
// http.Request and http.ResponseWriter.
type Ctx struct {
	context.Context
	Request  *http.Request
	Response http.ResponseWriter
}

// Param returns a URL path parameter by name.
// For a route registered as /users/{id}, use ctx.Param("id").
func (c *Ctx) Param(key string) string {
	return c.Request.PathValue(key)
}

// Query returns a URL query parameter by name.
func (c *Ctx) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// Header returns a request header value by name.
func (c *Ctx) Header(key string) string {
	return c.Request.Header.Get(key)
}

// SetHeader sets a response header.
func (c *Ctx) SetHeader(key, value string) {
	c.Response.Header().Set(key, value)
}
