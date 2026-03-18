// Package utils provides internal utility functions for the swagify package.
package utils

import (
	"strings"
	"unicode"
)

// NormalizePath converts framework-specific path parameters to OpenAPI format.
// For example: /users/:id => /users/{id}, /users/*filepath => /users/{filepath}
func NormalizePath(path string) string {
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if len(seg) == 0 {
			continue
		}
		// Fiber/Gin style :param
		if seg[0] == ':' {
			segments[i] = "{" + seg[1:] + "}"
		}
		// Wildcard *param
		if seg[0] == '*' {
			name := seg[1:]
			if name == "" {
				name = "wildcard"
			}
			segments[i] = "{" + name + "}"
		}
	}
	return strings.Join(segments, "/")
}

// ExtractPathParams returns the parameter names from a path.
// For example: /users/:id/posts/:postId => ["id", "postId"]
func ExtractPathParams(path string) []string {
	var params []string
	segments := strings.Split(path, "/")
	for _, seg := range segments {
		if len(seg) == 0 {
			continue
		}
		if seg[0] == ':' {
			params = append(params, seg[1:])
		}
		if seg[0] == '*' {
			name := seg[1:]
			if name == "" {
				name = "wildcard"
			}
			params = append(params, name)
		}
		// Also handle {param} style
		if seg[0] == '{' && seg[len(seg)-1] == '}' {
			params = append(params, seg[1:len(seg)-1])
		}
	}
	return params
}

// GenerateOperationID creates a stable operation ID from method and path.
// For example: POST /users/{id}/posts => postUsersIdPosts
func GenerateOperationID(method, path string) string {
	// Clean the path
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	path = strings.ReplaceAll(path, ":", "")

	segments := strings.Split(path, "/")
	var parts []string
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		parts = append(parts, capitalize(seg))
	}

	return strings.ToLower(method) + strings.Join(parts, "")
}

// capitalize capitalizes the first letter of a string.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// ToSnakeCase converts a CamelCase string to snake_case.
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToCamelCase converts a snake_case or kebab-case string to CamelCase.
func ToCamelCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var result strings.Builder
	for _, p := range parts {
		result.WriteString(capitalize(p))
	}
	return result.String()
}
