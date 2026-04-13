package swagify

import "reflect"

// Option configures a route's documentation metadata.
type Option func(*route)

// Summary sets the operation summary shown in the docs.
func Summary(s string) Option {
	return func(r *route) { r.summary = s }
}

// Description sets a longer operation description.
func Description(d string) Option {
	return func(r *route) { r.description = d }
}

// Tags groups the operation under one or more named sections in the docs.
func Tags(tags ...string) Option {
	return func(r *route) { r.tags = append(r.tags, tags...) }
}

// OperationID sets a custom unique operation identifier.
func OperationID(id string) Option {
	return func(r *route) { r.operationID = id }
}

// Deprecated marks the operation as deprecated in the docs.
func Deprecated() Option {
	return func(r *route) { r.deprecated = true }
}

// SuccessStatus overrides the default success HTTP status code.
func SuccessStatus(code int) Option {
	return func(r *route) { r.successStatus = code }
}

// Response documents an additional response for a specific status code.
func Response(status int, model any, description string) Option {
	return func(r *route) {
		if r.responses == nil {
			r.responses = make(map[int]routeResponse)
		}
		var t reflect.Type
		if model != nil {
			t = reflect.TypeOf(model)
		}
		r.responses[status] = routeResponse{
			statusCode:  status,
			description: description,
			typ:         t,
		}
	}
}

// ErrorResponse is an alias for Response, used to document error responses.
func ErrorResponse(status int, model any, description string) Option {
	return Response(status, model, description)
}

// Security sets per-route security requirements.
// Each entry is a map of scheme name to required scopes.
func Security(schemes ...map[string][]string) Option {
	return func(r *route) { r.security = append(r.security, schemes...) }
}

// SecurityBearer requires a bearer token for this route.
func SecurityBearer() Option {
	return Security(map[string][]string{"bearerAuth": {}})
}

// SecurityAPIKey requires an API key for this route.
func SecurityAPIKey() Option {
	return Security(map[string][]string{"apiKeyAuth": {}})
}

// SecurityBasic requires HTTP basic auth for this route.
func SecurityBasic() Option {
	return Security(map[string][]string{"basicAuth": {}})
}
