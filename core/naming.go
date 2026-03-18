package core

import (
	"strings"
	"unicode"
)

// ToSchemaName converts a Go type name to a clean schema component name.
func ToSchemaName(name string) string {
	if name == "" {
		return "Anonymous"
	}
	return name
}

// ToOperationID generates a stable operation ID from an HTTP method and path.
// Example: POST /users/{id}/posts => postUsersIdPosts
func ToOperationID(method, path string) string {
	cleaned := strings.NewReplacer("{", "", "}", "", ":", "").Replace(path)
	segments := strings.Split(cleaned, "/")
	var parts []string
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		parts = append(parts, capitalizeFirst(seg))
	}
	return strings.ToLower(method) + strings.Join(parts, "")
}

// ToJSONFieldName converts a Go field name to a JSON-style name (camelCase).
func ToJSONFieldName(name string) string {
	if name == "" {
		return name
	}
	runes := []rune(name)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
