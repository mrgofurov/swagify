package openapi

import "github.com/swagify/core"

// buildSecuritySchemes converts internal security scheme definitions
// to OpenAPI SecuritySchemeObjects.
func buildSecuritySchemes(schemes []core.SecurityScheme) map[string]*SecuritySchemeObject {
	if len(schemes) == 0 {
		return nil
	}

	result := make(map[string]*SecuritySchemeObject)
	for _, s := range schemes {
		name := inferSecuritySchemeName(s)
		result[name] = &SecuritySchemeObject{
			Type:         s.Type,
			Scheme:       s.Scheme,
			BearerFormat: s.BearerFormat,
			Name:         s.Name,
			In:           s.In,
			Description:  s.Description,
		}
	}
	return result
}

// inferSecuritySchemeName generates a name for a security scheme
// based on its type and configuration.
func inferSecuritySchemeName(s core.SecurityScheme) string {
	if s.Description != "" {
		// Try to use description-based naming
	}

	switch s.Type {
	case "http":
		if s.Scheme == "bearer" {
			return "bearerAuth"
		}
		if s.Scheme == "basic" {
			return "basicAuth"
		}
		return "httpAuth"
	case "apiKey":
		return "apiKeyAuth"
	case "oauth2":
		return "oauth2"
	case "openIdConnect":
		return "openIdConnect"
	default:
		return "auth"
	}
}

// BearerSecurity creates a bearer token security scheme.
func BearerSecurity(description ...string) core.SecurityScheme {
	desc := "Bearer token authentication"
	if len(description) > 0 {
		desc = description[0]
	}
	return core.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  desc,
	}
}

// APIKeySecurity creates an API key security scheme.
func APIKeySecurity(name, in string, description ...string) core.SecurityScheme {
	desc := "API key authentication"
	if len(description) > 0 {
		desc = description[0]
	}
	return core.SecurityScheme{
		Type:        "apiKey",
		Name:        name,
		In:          in,
		Description: desc,
	}
}

// BasicSecurity creates a basic auth security scheme.
func BasicSecurity(description ...string) core.SecurityScheme {
	desc := "Basic HTTP authentication"
	if len(description) > 0 {
		desc = description[0]
	}
	return core.SecurityScheme{
		Type:        "http",
		Scheme:      "basic",
		Description: desc,
	}
}
