// Package swagify provides automatic OpenAPI 3.1 documentation generation
// for Go web frameworks including Fiber and Gin.
//
// Swagify takes a code-first approach: define your types and handlers,
// and swagify generates a complete OpenAPI specification with a beautiful
// interactive docs UI — no comments or annotations required.
//
// Quick start with Fiber:
//
//	app := fiber.New()
//	api := swagify.NewFiber(app)
//
//	api.POST("/users", createUserHandler,
//	    swagify.Summary("Create a user"),
//	    swagify.Tags("Users"),
//	    swagify.WithRequest(CreateUserRequest{}),
//	    swagify.WithResponse(UserResponse{}),
//	)
//
//	api.RegisterOpenAPI()
//	api.RegisterDocs()
//	app.Listen(":8080")
//
// Quick start with Gin:
//
//	r := gin.Default()
//	api := swagify.NewGin(r)
//
//	api.POST("/users", createUserHandler,
//	    swagify.Summary("Create a user"),
//	    swagify.Tags("Users"),
//	)
//
//	api.RegisterOpenAPI()
//	api.RegisterDocs()
//	r.Run(":8080")
package swagify

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/swagify/core"
	"github.com/swagify/router"
)

// --- Framework Constructors ---

// NewFiber creates a new swagify adapter for the Fiber web framework.
func NewFiber(app *fiber.App, configs ...router.FiberConfig) *router.FiberAdapter {
	return router.NewFiber(app, configs...)
}

// NewGin creates a new swagify adapter for the Gin web framework.
func NewGin(engine *gin.Engine, configs ...router.GinConfig) *router.GinAdapter {
	return router.NewGin(engine, configs...)
}

// --- Route Options (re-exported for clean top-level API) ---

// RouteOption configures a route's documentation metadata.
type RouteOption = router.RouteOption

// Summary sets the operation summary.
func Summary(s string) RouteOption {
	return router.Summary(s)
}

// Description sets the operation description.
func Description(d string) RouteOption {
	return router.Description(d)
}

// Tags sets the operation tags for grouping.
func Tags(tags ...string) RouteOption {
	return router.Tags(tags...)
}

// OperationID sets a custom operation ID.
func OperationID(id string) RouteOption {
	return router.OperationID(id)
}

// Deprecated marks the operation as deprecated.
func Deprecated() RouteOption {
	return router.DeprecatedOp()
}

// SuccessStatus sets the default success status code.
func SuccessStatus(code int) RouteOption {
	return router.SuccessStatus(code)
}

// Response adds a custom response definition for a status code.
func Response(status int, model any, description string) RouteOption {
	return router.Response(status, model, description)
}

// ErrorResponse adds an error response definition.
func ErrorResponse(status int, model any, description string) RouteOption {
	return router.ErrorResponse(status, model, description)
}

// Security sets security requirements for a route.
func Security(schemes ...map[string][]string) RouteOption {
	return router.Security(schemes...)
}

// SecurityBearer adds bearer token security requirement to a route.
func SecurityBearer() RouteOption {
	return router.SecurityBearer()
}

// SecurityAPIKey adds API key security requirement to a route.
func SecurityAPIKey() RouteOption {
	return router.SecurityAPIKey()
}

// SecurityBasic adds basic auth security requirement to a route.
func SecurityBasic() RouteOption {
	return router.SecurityBasic()
}

// WithRequest sets the request body type for untyped handlers.
func WithRequest(model any) RouteOption {
	return router.WithRequest(model)
}

// WithResponse sets the response body type for untyped handlers.
func WithResponse(model any) RouteOption {
	return router.WithResponse(model)
}

// QueryParams sets the query parameters type for documentation.
func QueryParams(model any) RouteOption {
	return router.QueryParams(model)
}

// PathParams sets the path parameters type for documentation.
func PathParams(model any) RouteOption {
	return router.PathParams(model)
}

// HeaderParams sets the header parameters type for documentation.
func HeaderParams(model any) RouteOption {
	return router.HeaderParams(model)
}

// RequestContentType overrides the request content type.
func RequestContentType(ct string) RouteOption {
	return router.RequestContentType(ct)
}

// ResponseContentType overrides the response content type.
func ResponseContentType(ct string) RouteOption {
	return router.ResponseContentType(ct)
}

// --- Typed Handler Registration (Fiber) ---

// POST registers a typed POST handler for Fiber with automatic schema inference.
func POST[Req any, Res any](f *router.FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	router.TypedPOST[Req, Res](f, path, handler, opts...)
}

// GET registers a typed GET handler for Fiber with automatic schema inference.
func GET[Req any, Res any](f *router.FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	router.TypedGET[Req, Res](f, path, handler, opts...)
}

// PUT registers a typed PUT handler for Fiber with automatic schema inference.
func PUT[Req any, Res any](f *router.FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	router.TypedPUT[Req, Res](f, path, handler, opts...)
}

// PATCH registers a typed PATCH handler for Fiber with automatic schema inference.
func PATCH[Req any, Res any](f *router.FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	router.TypedPATCH[Req, Res](f, path, handler, opts...)
}

// DELETE registers a typed DELETE handler for Fiber with automatic schema inference.
func DELETE[Req any, Res any](f *router.FiberAdapter, path string, handler func(*fiber.Ctx, Req) (Res, error), opts ...RouteOption) {
	router.TypedDELETE[Req, Res](f, path, handler, opts...)
}

// --- Core Types (re-exported for convenience) ---

// Info represents API information metadata.
type Info = core.Info

// Contact represents API contact information.
type Contact = core.Contact

// License represents API license information.
type License = core.License

// Server represents an API server.
type Server = core.Server

// SecurityScheme represents a security scheme definition.
type SecurityScheme = core.SecurityScheme

// Tag represents an API tag for grouping operations.
type Tag = core.Tag

// --- Config Types (re-exported) ---

// FiberConfig holds configuration for the Fiber adapter.
type FiberConfig = router.FiberConfig

// GinConfig holds configuration for the Gin adapter.
type GinConfig = router.GinConfig
