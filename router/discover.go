package router

import (
	"strings"
	"unicode"
)

// DiscoverOptions configures how route auto-discovery behaves.
type DiscoverOptions struct {
	// ExcludePaths excludes routes whose paths start with any of these prefixes.
	// Example: []string{"/health", "/metrics"}
	ExcludePaths []string

	// IncludePaths limits discovery to routes whose paths start with one of these prefixes.
	// If empty, all paths are included (except excluded ones).
	IncludePaths []string

	// ExcludeInternal skips internal framework routes (docs, openapi, favicon, etc.).
	// Default: true
	ExcludeInternal *bool

	// AutoTags automatically generates tags from the first meaningful path segment.
	// /users/:id → tag "Users", /api/v1/orders → tag "Orders"
	// Default: true
	AutoTags *bool

	// AutoSummary automatically generates human-readable summaries from method + path.
	// GET /users/:id → "Get user by id"
	// Default: true
	AutoSummary *bool
}

// defaultInternalPaths are framework paths skipped when ExcludeInternal is true.
var defaultInternalPaths = []string{
	"/openapi.json",
	"/openapi.yaml",
	"/docs",
	"/swagger",
	"/favicon.ico",
	"/health",
	"/healthz",
	"/ready",
	"/readyz",
	"/livez",
}

// shouldInclude checks if a path should be included based on DiscoverOptions.
func shouldInclude(path string, opts DiscoverOptions) bool {
	excludeInternal := true
	if opts.ExcludeInternal != nil {
		excludeInternal = *opts.ExcludeInternal
	}

	// Check internal paths
	if excludeInternal {
		for _, internal := range defaultInternalPaths {
			if path == internal || strings.HasPrefix(path, internal+"/") {
				return false
			}
		}
	}

	// Check excluded paths
	for _, exc := range opts.ExcludePaths {
		if strings.HasPrefix(path, exc) {
			return false
		}
	}

	// Check included paths (if specified)
	if len(opts.IncludePaths) > 0 {
		included := false
		for _, inc := range opts.IncludePaths {
			if strings.HasPrefix(path, inc) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	return true
}

// autoTag extracts a tag name from the path by finding the first meaningful segment.
// /users/:id → "Users"
// /api/v1/orders/:id/items → "Orders"
func autoTag(path string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")

	for _, seg := range segments {
		if seg == "" {
			continue
		}
		// Skip version segments like v1, v2, api
		lower := strings.ToLower(seg)
		if lower == "api" || (len(lower) >= 2 && lower[0] == 'v' && isDigit(lower[1])) {
			continue
		}
		// Skip path parameters
		if seg[0] == ':' || seg[0] == '*' || (seg[0] == '{' && seg[len(seg)-1] == '}') {
			continue
		}
		// Capitalize and return
		return capitalize(lower)
	}

	return "Default"
}

// autoSummary generates a human-readable summary from method and path.
// GET /users → "List users"
// GET /users/:id → "Get user by id"
// POST /users → "Create user"
// PUT /users/:id → "Update user by id"
// PATCH /users/:id → "Partially update user by id"
// DELETE /users/:id → "Delete user by id"
func autoSummary(method, path string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")

	// Find meaningful segments (skip api, version prefixes, params)
	var resource string
	var params []string

	for _, seg := range segments {
		if seg == "" {
			continue
		}
		lower := strings.ToLower(seg)

		// Skip version segments
		if lower == "api" || (len(lower) >= 2 && lower[0] == 'v' && isDigit(lower[1])) {
			continue
		}

		// Collect param names
		if seg[0] == ':' {
			params = append(params, seg[1:])
			continue
		}
		if seg[0] == '*' {
			name := seg[1:]
			if name == "" {
				name = "path"
			}
			params = append(params, name)
			continue
		}
		if seg[0] == '{' && seg[len(seg)-1] == '}' {
			params = append(params, seg[1:len(seg)-1])
			continue
		}

		// Use last non-param segment as the resource
		resource = lower
	}

	if resource == "" {
		resource = "resource"
	}

	// Singularize resource for single-item operations
	singular := singularize(resource)

	// Build the summary
	var action string
	hasParams := len(params) > 0

	switch method {
	case "GET":
		if hasParams {
			action = "Get " + singular
		} else {
			action = "List " + resource
		}
	case "POST":
		action = "Create " + singular
	case "PUT":
		if hasParams {
			action = "Update " + singular
		} else {
			action = "Replace " + resource
		}
	case "PATCH":
		if hasParams {
			action = "Partially update " + singular
		} else {
			action = "Update " + resource
		}
	case "DELETE":
		if hasParams {
			action = "Delete " + singular
		} else {
			action = "Delete " + resource
		}
	case "HEAD":
		action = "Check " + singular
	case "OPTIONS":
		action = "Options for " + resource
	default:
		action = strings.ToLower(method) + " " + resource
	}

	// Append parameter context
	if hasParams {
		action += " by " + strings.Join(params, " and ")
	}

	// Capitalize first letter
	return capitalizeFirst(action)
}

// singularize does a very basic English singularization.
func singularize(word string) string {
	if word == "" {
		return word
	}
	if strings.HasSuffix(word, "ies") && len(word) > 3 {
		return word[:len(word)-3] + "y"
	}
	if strings.HasSuffix(word, "ses") || strings.HasSuffix(word, "xes") || strings.HasSuffix(word, "zes") {
		return word[:len(word)-2]
	}
	if strings.HasSuffix(word, "s") && !strings.HasSuffix(word, "ss") {
		return word[:len(word)-1]
	}
	return word
}

// capitalize capitalizes the first letter of a word.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// capitalizeFirst capitalizes only the very first letter of a sentence.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// isDigit checks if a byte is a digit.
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// boolPtr is a helper to create a pointer to a bool.
func boolPtr(v bool) *bool {
	return &v
}
