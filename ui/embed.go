// Package ui provides the embedded Swagger UI and handler registration.
package ui

import (
	"html/template"
	"net/http"
	"strings"

	_ "embed"
)

//go:embed swagger/index.html
var swaggerHTML string

// templateData holds the values injected into the UI template.
type templateData struct {
	Title   string
	SpecURL string
}

// renderHTML renders the Swagger UI HTML with the given spec URL and title.
func renderHTML(specURL string, title ...string) (string, error) {
	t := "API Documentation"
	if len(title) > 0 && title[0] != "" {
		t = title[0]
	}

	tmpl, err := template.New("docs").Parse(swaggerHTML)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, templateData{Title: t, SpecURL: specURL}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RegisterHTTP mounts the Swagger UI on a standard http.ServeMux.
func RegisterHTTP(mux *http.ServeMux, docsPath, specURL string) {
	html, err := renderHTML(specURL)
	if err != nil {
		panic("swagify: failed to render docs UI: " + err.Error())
	}
	mux.HandleFunc("GET "+docsPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html)) //nolint:errcheck
	})
}

// RegisterHTTPWithAuth mounts the Swagger UI behind a middleware.
func RegisterHTTPWithAuth(mux *http.ServeMux, docsPath, specURL string, mw func(http.Handler) http.Handler) {
	html, err := renderHTML(specURL)
	if err != nil {
		panic("swagify: failed to render docs UI: " + err.Error())
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html)) //nolint:errcheck
	})
	mux.Handle("GET "+docsPath, mw(handler))
}
