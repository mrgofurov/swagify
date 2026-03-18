// Package errors provides structured error types for the swagify package.
package errors

import "fmt"

// SchemaError is returned when schema generation fails for a type.
type SchemaError struct {
	TypeName string
	Message  string
}

func (e *SchemaError) Error() string {
	return fmt.Sprintf("swagify: schema error for type %q: %s", e.TypeName, e.Message)
}

// RegistrationError is returned when route registration fails.
type RegistrationError struct {
	Path    string
	Method  string
	Message string
}

func (e *RegistrationError) Error() string {
	return fmt.Sprintf("swagify: registration error for %s %s: %s", e.Method, e.Path, e.Message)
}

// HandlerError is returned when a handler has an unsupported signature.
type HandlerError struct {
	HandlerName string
	Message     string
}

func (e *HandlerError) Error() string {
	return fmt.Sprintf("swagify: handler error for %q: %s", e.HandlerName, e.Message)
}
