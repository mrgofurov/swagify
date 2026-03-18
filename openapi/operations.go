package openapi

import "fmt"

// hasRequestBody returns true if the HTTP method typically has a request body.
func hasRequestBody(method string) bool {
	switch method {
	case "POST", "PUT", "PATCH":
		return true
	default:
		return false
	}
}

// defaultStatusDescription returns a human-readable description for an HTTP status code.
func defaultStatusDescription(status int) string {
	descriptions := map[int]string{
		200: "Successful response",
		201: "Created successfully",
		204: "No content",
		400: "Bad request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not found",
		409: "Conflict",
		422: "Validation error",
		500: "Internal server error",
	}
	if desc, ok := descriptions[status]; ok {
		return desc
	}
	return fmt.Sprintf("Response %d", status)
}
