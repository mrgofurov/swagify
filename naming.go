package swagify

import (
	"strings"
	"unicode"
)

// autoTag extracts a tag from the first meaningful path segment.
//
//	/users/{id}       → "Users"
//	/api/v1/orders    → "Orders"
func autoTag(path string) string {
	for _, seg := range strings.Split(strings.Trim(path, "/"), "/") {
		if seg == "" {
			continue
		}
		lower := strings.ToLower(seg)
		if lower == "api" || (len(lower) >= 2 && lower[0] == 'v' && lower[1] >= '0' && lower[1] <= '9') {
			continue
		}
		if seg[0] == ':' || seg[0] == '*' || (seg[0] == '{' && seg[len(seg)-1] == '}') {
			continue
		}
		return capitalizeFirst(lower)
	}
	return ""
}

// autoSummary generates a human-readable summary from the HTTP method and path.
//
//	GET  /users        → "List users"
//	GET  /users/{id}   → "Get user by id"
//	POST /users        → "Create user"
func autoSummary(method, path string) string {
	var resource string
	var params []string

	for _, seg := range strings.Split(strings.Trim(path, "/"), "/") {
		if seg == "" {
			continue
		}
		lower := strings.ToLower(seg)

		if lower == "api" || (len(lower) >= 2 && lower[0] == 'v' && lower[1] >= '0' && lower[1] <= '9') {
			continue
		}
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
		resource = lower
	}

	if resource == "" {
		resource = "resource"
	}

	singular := singularize(resource)
	hasParams := len(params) > 0

	var action string
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
	default:
		action = strings.ToLower(method) + " " + resource
	}

	if hasParams {
		action += " by " + strings.Join(params, " and ")
	}
	return capitalizeFirst(action)
}

func singularize(word string) string {
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

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
