package router

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// RouteOption is a function that modifies a route's metadata during registration.
// This enables a clean, composable API for enriching route documentation.
type RouteOption func(*Route)

// Summary sets the operation summary.
func Summary(s string) RouteOption {
	return func(r *Route) {
		r.Summary = s
	}
}

// Description sets the operation description.
func Description(d string) RouteOption {
	return func(r *Route) {
		r.Description = d
	}
}

// Tags sets the operation tags.
func Tags(tags ...string) RouteOption {
	return func(r *Route) {
		r.Tags = append(r.Tags, tags...)
	}
}

// OperationID sets a custom operation ID.
func OperationID(id string) RouteOption {
	return func(r *Route) {
		r.OperationID = id
	}
}

// DeprecatedOp marks the operation as deprecated.
func DeprecatedOp() RouteOption {
	return func(r *Route) {
		r.Deprecated = true
	}
}

// SuccessStatus sets the default success status code.
func SuccessStatus(code int) RouteOption {
	return func(r *Route) {
		r.SuccessStatus = code
	}
}

// Response adds a custom response definition for a status code.
func Response(status int, model any, description string) RouteOption {
	return func(r *Route) {
		if r.Responses == nil {
			r.Responses = make(map[int]RouteResponse)
		}
		var t reflect.Type
		if model != nil {
			t = reflect.TypeOf(model)
		}
		r.Responses[status] = RouteResponse{
			StatusCode:  status,
			Description: description,
			Type:        t,
		}
	}
}

// ErrorResponse adds an error response definition.
func ErrorResponse(status int, model any, description string) RouteOption {
	return Response(status, model, description)
}

// Security sets the security requirements for this route.
func Security(schemes ...map[string][]string) RouteOption {
	return func(r *Route) {
		r.Security = append(r.Security, schemes...)
	}
}

// SecurityBearer adds a bearer token security requirement.
func SecurityBearer() RouteOption {
	return func(r *Route) {
		r.Security = append(r.Security, map[string][]string{
			"bearerAuth": {},
		})
	}
}

// SecurityAPIKey adds an API key security requirement.
func SecurityAPIKey() RouteOption {
	return func(r *Route) {
		r.Security = append(r.Security, map[string][]string{
			"apiKeyAuth": {},
		})
	}
}

// SecurityBasic adds a basic auth security requirement.
func SecurityBasic() RouteOption {
	return func(r *Route) {
		r.Security = append(r.Security, map[string][]string{
			"basicAuth": {},
		})
	}
}

// FiberMiddleware adds a middleware to the route.
func FiberMiddleware(middleware ...fiber.Handler) RouteOption {
	return func(r *Route) {
		r.FiberMiddlewares = append(r.FiberMiddlewares, middleware...)
	}
}

// GinMiddleware adds a middleware to the route.
func GinMiddleware(middleware ...gin.HandlerFunc) RouteOption {
	return func(r *Route) {
		r.GinMiddlewares = append(r.GinMiddlewares, middleware...)
	}
}

// QueryParams sets the query parameters type for documentation.
func QueryParams(model any) RouteOption {
	return func(r *Route) {
		if model != nil {
			r.QueryType = reflect.TypeOf(model)
		}
	}
}

// PathParams sets the path parameters type for documentation.
func PathParams(model any) RouteOption {
	return func(r *Route) {
		if model != nil {
			r.PathType = reflect.TypeOf(model)
		}
	}
}

// HeaderParams sets the header parameters type for documentation.
func HeaderParams(model any) RouteOption {
	return func(r *Route) {
		if model != nil {
			r.HeaderType = reflect.TypeOf(model)
		}
	}
}

// RequestContentType overrides the request content type.
func RequestContentType(ct string) RouteOption {
	return func(r *Route) {
		r.RequestContentType = ct
	}
}

// ResponseContentType overrides the response content type.
func ResponseContentType(ct string) RouteOption {
	return func(r *Route) {
		r.ResponseContentType = ct
	}
}

// applyOptions applies a list of RouteOptions to a route.
func applyOptions(route *Route, opts []RouteOption) {
	for _, opt := range opts {
		opt(route)
	}
}
